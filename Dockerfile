FROM postgres
LABEL authors="slavaruswarrior"

ENTRYPOINT ["top", "-b"]