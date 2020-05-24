from flask import Flask
from flask import make_response
from flask import request
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
model._make_predict_function()

logger.info("Loaded model from disk")

app = Flask(__name__)
app.config['UPLOAD_FOLDER'] = "./"

def predict(img_filename):
    global model

    img = skimage.io.imread(img_filename)
    img = skimage.transform.resize(img, (224, 224))
    img = keras.applications.mobilenet.preprocess_input(img)
    img = np.expand_dims(img, 0)

    predictions = model.predict(img).reshape(-1).tolist()
    top_pred = np.argmax(predictions)
    disease_name = class_list[top_pred]

    return predictions, class_list


@app.route('/dermatologia', methods=['POST'])
def Compute():
    img = request.files['file']
    img.save('image.tmp')

    predictions, class_list = predict("image.tmp")
    toReturn = {"Types": class_list, "Probabilities": predictions}
    return make_response(toReturn, 200)


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8000)
