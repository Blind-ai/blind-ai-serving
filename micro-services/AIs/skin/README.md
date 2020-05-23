## Build and run with docker

<h2>Without GPU<h2>

`docker build . -t blind/ai/skin`

`docker run -p 8501:8501 --rm --name blind_preprocess blind/ai/skin`

return json:
{'Types': [], 'Probabilities': []}
