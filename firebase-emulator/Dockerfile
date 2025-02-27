FROM node:22-alpine

WORKDIR /firebase

RUN npm install -g firebase-tools

COPY firebase.json .

CMD ["firebase", "emulators:start", "--project", "demo-backend"]
