# How to setup autocompletion for APL editing in vscode

1. Install protoc-gen-jsonschema

```
go install github.com/chrusty/protoc-gen-jsonschema/cmd/protoc-gen-jsonschema@latest
```

2. Generate Schema

```
mkdir jsonschema

protoc \
--jsonschema_out=./jsonschema \
--proto_path=./proto/ \
--jsonschema_opt=json_fieldnames \
--jsonschema_opt=enforce_oneof \
--jsonschema_opt=disallow_additional_properties \
./proto/apl.proto
```

3. Edit APL files

The project settings file for vscode should now pick that up (.vscode/settings.json) should point the schema "jsonschema" directory. The extension ".apl.json" should activate the generated schema.

Create an APL file named "myrotation.apl.json" replacing `myrotation` with whatever name you prefer. 

An empty rotation looks something like 
```
{
    "type": "TypeAPL",
    "prepullActions": [],
    "priorityList": []
}
```

You can then add prepull actions and the main priority list.

Here is what a example prepull action looks like. It casts spell 1 at -1 seconds.
```
        {
            "action": {
                "castSpell": {
                    "spellId": {
                        "spellId": 1
                    }
                }
            },
            "doAt": "-1s"
        }
```

Here is an example action from the priorityList of elemental shaman. This is the check to see if flameshock dot is applied before casting lavaburst.

Condition is checking that "cmp" (compare) that the "lhs" (left hand side) is "OpGt" (greater than) the "rhs" (right hand side). 

In this case it is checking that dot remaining time for flameshock is greater than the cast time of lava burst.

```
        {
            "action": {
                "condition": {
                    "cmp": {
                        "op": "OpGt",
                        "lhs": {
                            "dotRemainingTime": {
                                "spellId": {
                                    "spellId": 49233
                                }
                            }
                        },
                        "rhs": {
                            "spellCastTime": {
                                "spellId": {
                                    "spellId": 60043
                                }
                            }
                        }
                    }
                },
                "castSpell": {
                    "spellId": {
                        "spellId": 60043
                    }
                }
            }
        },
```

5. Insert rotation into JSON file for import.

You export your current settings in the sim (Export->JSON). Save the export as a file. Replace the `"rotation": {}` part of the export with your custom json rotation. (Just replace the `{}` leaving the `"rotation":` )

In the sim click (Import->JSON) and choose your edited JSON file, your rotation should appear!