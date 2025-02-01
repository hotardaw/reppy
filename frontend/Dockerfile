FROM node:latest

WORKDIR /app

COPY package*.json ./

RUN npm install
RUN npm install --save-dev @types/react @types/react-dom typescript @types/node

COPY . .

EXPOSE 8080

CMD ["npm", "run", "dev"]