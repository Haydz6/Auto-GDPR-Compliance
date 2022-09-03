# Auto GDPR Compliance

Automatically scans through your messages and fulfills GDPR compliance messages.

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