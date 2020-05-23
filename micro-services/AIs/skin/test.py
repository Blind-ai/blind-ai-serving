import numpy as np
import keras
import skimage.io
import skimage.transform
from keras.models import model_from_json
from loguru import logger

class_list = ['mel', 'bkl', 'bcc', 'akiec', 'vasc', 'df']

logger.info("Loading model from disk...")

json_file = open('model_definitive.json', 'r')
loaded_model_json = json_file.read()
json_file.close()
model = model_from_json(loaded_model_json)
model.load_weights("model_definitive.h5")

logger.info("Loaded model from disk")

def predict(img_filename):
    img = skimage.io.imread(img_filename)
    img = skimage.transform.resize(img, (224, 224))
    img = keras.applications.mobilenet.preprocess_input(img)
    img = np.expand_dims(img, 0)

    predictions     = model.predict(img).reshape(-1)
    top_pred        = np.argmax(predictions)
    disease_name    = class_list[top_pred]

    print(type(predictions))
    return predictions, disease_name


prediction, disease_name = predict("image.tmp")
print(prediction, disease_name)