# Auto GDPR Compliance

Automatically scans through your messages and fulfills GDPR compliance messages.

# NOTICE
Enabling DeleteGDPRMessagesAfterFulfilled will permanently delete the GDPR messages forever.
Archiving system messages will not put them into the archive inbox. (Thanks roblox 😡)

# How to use

## Inserting ROBLOSECURITY
You have to insert your ROBLOSECURITY into the relevant text file as such:

```
Put your ROBLOSECURITY below!!!
REPLACETHISWITHYOURROBLOSECURITY
```

## Adding keys
You must insert the following into the DataKeys.json file

```
{
    "PLACEIDHERE": {
        "DATASTORETYPE": {
            "DATASTORENAMEHERE": {
                "SCOPEHERE": [
                    "KEY%USERID"
                ]
            }
        }
    }
}
```

Example: 
```
{
    "10792426231": {
        "DataStore": {
            "PlayerData": {
                "": [
                    "Pets_%USERID",
                    "Money_%USERID"
                ]
            }
        },
        "OrderedDataStore": {
            "Money": {
                "": [
                    "%USERID"
                ]
            }
        }
    }
}
```