FROM postgres:13-bullseye
RUN localedef -i fi_FI -c -f UTF-8 -A /usr/share/locale/locale.alias fi_FI.UTF-8
ENV LANG fi_FI.utf8
ENV POSTGRES_USER tester
ENV POSTGRES_PASSWORD testing
ENV POSTGRES_DB todos