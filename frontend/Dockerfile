FROM node:latest

WORKDIR /app

COPY package*.json ./

RUN npm clean-install

COPY . .

ENV VITE_HOST=0.0.0.0

EXPOSE 8080

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]