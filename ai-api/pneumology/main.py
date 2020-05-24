import skimage
import torch
import torchxrayvision as xrv
import torchvision
import os

from flask import Flask
from flask import request
from flask import make_response
from loguru import logger

logger.info("Initializing the model. This may take longer if it's the first start.")
model = xrv.models.DenseNet(weights="all")

app = Flask(__name__)
app.config['UPLOAD_FOLDER'] = './'
cuda_enabled=False

def camelCase(string):
    string = string.replace('_', ' ').replace('-', ' ')
    words = string.split(' ')
    formatted = []
    for i in range(0, len(words)):
        word = words[i]
        if not word.isalnum():
            continue

        if i == 0:
            formatted.append(word.lower())
        else:
            formatted.append(word.capitalize())
    return ''.join(formatted)
        

    

@app.route('/pneumology', methods=['POST'])
def Compute():
    global model
    global cuda_enabled

    img = request.files['file']
    img.save("image.tmp")

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

        result = {}
        output["preds"] = dict(zip(xrv.datasets.default_pathologies, preds[0].detach().numpy()))
        predictions = output["preds"].items()
        for entry in predictions:
            name = camelCase(entry[0])
            prediction = float(entry[1])
            result[name] = prediction

        return make_response(result, 200)


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8000)
