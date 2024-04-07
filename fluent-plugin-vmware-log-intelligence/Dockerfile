FROM ruby:2.5
MAINTAINER slk@vmware.com

RUN apt-get update && apt-get install -y \
  build-essential \
  nodejs

RUN mkdir -p /app

COPY ./ /app/
WORKDIR /app
RUN gem build fluent-plugin-vmware-log-intelligence.gemspec

RUN gem install bundler && bundle install --jobs 20 --retry 5
RUN gem install fluent-plugin-vmware-log-intelligence
RUN gem list

WORKDIR /
RUN rm -rf /app