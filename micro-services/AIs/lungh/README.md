## Build and run with docker

<h2>Without GPU<h2>

`docker build . -t blind/ai/lungh`

`docker run -p 8501:8501 --rm --name blind_preprocess blind/ai/lungh`

<h2>With GPU<h2>

`docker build . -t blind/ai/lungh`

`docker run -p 8501:8501 --rm --name blind_preprocess blind/ai/lungh`
