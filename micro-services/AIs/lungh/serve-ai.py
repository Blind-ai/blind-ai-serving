import skimage
import torch
import torchxrayvision as xrv
import torchvision
import os

from flask import Flask
from flask import request
from flask import make_response
from loguru import logger


if os.environ['GPU'] == "false":
    cuda_enabled = False
    print("gpu disabled")
else:
    cuda_enabled = True
    print("gpu enabled")

logger.info("Initializing the model. This may take longer if it's the first start.")
model = xrv.models.DenseNet(weights="all")

app = Flask(__name__)

@app.route('/api/lungh/compute', methods=['POST'])
def Compute():
    global model
    global cuda_enabled
    print(request.args)
    with open("image.tmp", "wb") as w:
        w.write(request.data)
        w.close()

    img = skimage.io.imread("image.tmp")
    img = xrv.datasets.normalize(img, 255)

    if len(img.shape) > 2:
        img = img[:, :, 0]
    if len(img.shape) < 2:
        print("error, dimension lower than 2 for image")

    # Add color channel
    img = img[None, :, :]

    transform = torchvision.transforms.Compose([xrv.datasets.XRayCenterCrop(),
                                                xrv.datasets.XRayResizer(224)])

    img = transform(img)

    output = {}
    with torch.no_grad():
        img = torch.from_numpy(img).unsqueeze(0)
        if cuda_enabled:
            img = img.cuda()
            model = model.cuda()
        preds = model(img).cpu()

        output["preds"] = dict(zip(xrv.datasets.default_pathologies, preds[0].detach().numpy()))
        sortedPreds = sorted(output["preds"].items(), key=lambda x: x[1], reverse=True)[0]
        toReturn = {"Type": sortedPreds[0], "Probability": float(sortedPreds[1])}
        return make_response(toReturn, 200)


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8501)