import skimage
import torch
import torch.nn.functional as F
import torchvision, torchvision.transforms
import torchxrayvision as xrv
from loguru import logger

cuda_enabled = True

logger.info("Initializing the model. This may take longer if it's the first start.")
model = xrv.models.DenseNet(weights="all")

img = skimage.io.imread("test.png")
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
    output["preds"] = dict(zip(xrv.datasets.default_pathologies,preds[0].detach().numpy()))
    print(output)