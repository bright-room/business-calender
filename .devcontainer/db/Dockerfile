FROM postgres:latest

ARG DB_LANG=ja_JP

RUN apt-get update -y && \
    apt-get install -y --no-install-recommends locales &&  \
    rm -rf /var/lib/apt/lists/* \ &&  \
    localedef -i $DB_LANG -c -f UTF-8 -A /usr/share/locale/locale.alias $DB_LANG.utf8

ENV LANG $DB_LANG.utf8
