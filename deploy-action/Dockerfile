FROM debian:9.5-slim

LABEL "com.github.actions.name"="Deploy"
LABEL "com.github.actions.description"="Rebuild docker containers command"
LABEL "com.github.actions.icon"="mic"
LABEL "com.github.actions.color"="purple"

LABEL "homepage"="http://github.com/actions"
LABEL "maintainer"="Kaspars <smilga.kaspars@github.com>"

ADD deploy.sh /deploy.sh
ENTRYPOINT ["/deploy.sh"]
