FROM node:5

ENV GULP_INSTALLS gulp gulp-filter gulp-main-bower-files gulp-uglify

RUN npm install -g gulp@latest && \
    npm install -D $GULP_INSTALLS

ENTRYPOINT ["gulp"]
