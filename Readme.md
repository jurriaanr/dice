**After cloning this repo, rename the folder to legion!**

To deploy the function exec (after installing gcloud tools: https://cloud.google.com/sdk/install) from the legiondice folder 

    gcloud functions deploy dice --entry-point RollDice --runtime go111 --trigger-http
    
To test locally exec

    go run legion/test

To update gcloud

    gcloud components update