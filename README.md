
# Instabot FB - Wit.ai enabled FB bot boilerplate

[![Build Status](https://travis-ci.org/yeyus/instabot-fb.svg?branch=master)](https://travis-ci.org/yeyus/instabot-fb)

## Introduction

This repository is a very simple Go boilerplate project integrating wit.ai and Facebook Messenger API to provide support for bot integration. It also supports fetching valid CA signed certificates through Let's Encrypt.

## Usage

Clone this repository to your server

```
git clone https://github.com/yeyus/instabot-fb.git
```

Create a ```.env``` file by following this format. All the information you need to complete this file is available by following [Setup wit.ai project](#setup-witai-project) and [Setup FB Page & Bot](#setup-fb-page-bot).

```
GENERATE_CERTS={TRUE|FALSE}
PORT={1-65535}
DOMAIN={FQDN from which you will serve the bot}
FB_ACCESS_TOKEN=...
FB_VERIFY_TOKEN={you choose this value here and set in FB Webhook config}
FB_SECRET_TOKEN=...
WITAI_TOKEN=...
```

**GENERATE_CERTS** choose TRUE if you want the bot to fetch certificates from Let's Encrypt, by using this mode you will incur into some limitations. PORT has to be 443 and DOMAIN should exactly match the FQDN from which you will serve the bot. This is the most straightforward mode if you don't want to setup a reverse proxy using Nginx (explained [here](#serving-through-nginx-reverse-proxy)). If you choose FALSE no certificates will be provisioned and the bot will server regular HTTP traffic in the selected PORT.

**PORT** choose a port number where the bot should serve traffic.

**DOMAIN** set the FQDN with no leading http:// or https://. It's very important if you choose the self provisioned certificate mode as a mismatched domain will certainly break the certificate validity.

**FB_ACCESS_TOKEN** this token is obtained by following [Setup FB Page & Bot](#setup-fb-page-bot) guide.

**FB_VERIFY_TOKEN** you have to set this value, you can either generate some random string or set to a secret string you would like. After you've chosen it you must set it into FB Messenger's Webhook dialog.

**FB_SECRET_TOKEN** this token is obtained by following [Setup FB Page & Bot](#setup-fb-page-bot) guide.

**WITAI_TOKEN** this token is obtained by following [Setup wit.ai project](#setup-witai-project) guide.

## Setup wit.ai project

Wit.ai is a platform were you can create and train a bot by using their easy graphical tools. That bot implementation can easily be connected to you bot logic contained in this project by implementing all actions specified in the *Stories* part of your bot.

In this boilerplate we will leverage their NLP and NLU capabilities in order to greatly simplify the implementation of our simple weather forecast bot.

You can visit the sample project implemented in this boilerplate [here](https://wit.ai/yeyus/WeatherApp) (you will need a wit.ai account for accessing the project) or follow their [Quickstart guide](https://wit.ai/docs/quickstart) in order to create a simple bot.

To obtain your **WITAI_TOKEN** going to the _Settings tab_ and copying the contents of _Server Access Token_.

## Setup FB Page & Bot

Wit.ai API is platform agnostic and will play well with most of other open chat platforms, we have chosen Facebook Messenger for this boilerplate but with some changes in the code you might make it work with you own.

Facebook Messenger's API works by means of using *webhooks*, their servers will reach yours in order to deliver user messages. Users will communicate with your bot using a Facebook Page.

You can follow this [tutorial](https://developers.facebook.com/docs/apps/register) in order to create an App, a Page and and register the Messenger platform on it.

## Serving directly with Let's Encrypt

Make sure that your ```.env``` file's **GENERATE_CERTS** field is set to **TRUE** and your **PORT** is set to 443. Your server should be pointed to by a domain name for Let's Encrypt to resolve and generate valid certificates for the endpoint.

## Serving through nginx reverse proxy

A more flexible approach is to serve regular HTTP and let some reverse proxy like Nginx to provision the certficates and handle all appropriate redirections. In this case you can use whatever **DOMAIN** and **PORT** you wish. 

You can follow [this tutorial](https://www.digitalocean.com/community/tutorials/how-to-secure-nginx-with-let-s-encrypt-on-ubuntu-16-04) in order to install Nginx with Let's Encrypt support.

Once you have completed it you can leverage reverse proxy feature to redirect the traffic to the bot app.

**/etc/nginx/sites-available/default**
```
# Instabot FB reverse proxy
server {
       server_name ssl1.example.com;
       listen 443 ssl http2;
       listen [::]:443 ssl http2;
       include snippets/ssl-params.conf;
       include snippets/ssl-ssl1.example.com.conf;

       location / {
		   proxy_pass http://127.0.0.1:10443;
       }
}
```

## Contribute

You can contribute to this project by opening Issues or sending PRs to correct or extend features.

## Thanks

The following projects where used in order to bring you this boilerplate:

* [go-messenger-bot](https://github.com/abhinavdahiya/go-messenger-bot) by Abhinav Dahiya
* [witgo](https://github.com/kurrik/witgo) by Arne Roomann-Kurrik
* and the Go team
