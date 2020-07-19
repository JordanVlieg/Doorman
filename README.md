# Doorman

A simple golang application that uses Twilio to automatically buzz open your apartment door for you.

Settings are configured in the .env file.  This can be deployed to any hosting provider through docker.

Ensure you configure your twilio account to submit a GET request to the `/knock` endpoint.

I have configured this to play the same buzz code tones repeatedly, this may not work in your case, edit the code in main.go