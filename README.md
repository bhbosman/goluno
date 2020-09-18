# goluno


This repo will connect to the Luno streaming service, using the Luno streaming protocol https://www.luno.com/en/developers/api#tag/Streaming-API-(beta)

For the application to connect, you need a valid luno streaming service account. the information must be stored in  ~/.luno/keys.json. 

The json file should look like this:
{
  "key": "{some key}",
  "secret": "{some secret}"
}

# What does it do?

## Connect to crypto currencies

The application will connect to the following currencies out of the box:
* XBTZAR
* XBTEUR
* XBTUGX
* XBTZMW
* ETHXBT
* BCHXBT

## Extracting data
### Text Listener
A Text Listener port will be opened on 127.0.0.1:3000. This port will publish updates in json format of the crypto currencies as they change in real time. Not all changes will propagate,  only changes in the top two lines of the bids and asks will 

### Protobuf derived protocol
This port, on 127.0.0.1:3001, is a protobuf derived stack stack implementation, I am currently working on. More on this later



