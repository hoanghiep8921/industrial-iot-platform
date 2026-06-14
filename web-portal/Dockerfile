FROM nginx:alpine
RUN apk add --no-cache tzdata
ENV TZ=Asia/Ho_Chi_Minh
COPY dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
