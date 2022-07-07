
### Update Tool 

## sample problem 

so how do we add this variable to assets.json  

```
{

         "absolute_power_direction_flag":
            {
               "name": "Charging when Enabled",
                "ui_type": "control",
                "type": "enum_slider",
                "var_type": "Bool",
                "value": false,
                "remote_enabled": true,
                "options":
                [
                   {
                    "name": "Charge",
                    "value": true
                    },
                    {
                    "name": "Discharge",
                    "value": false
                    }
                ]
           }
        }
```

At this location in asstes.json 

```
{
     "ess":
      {
      "asset_instances":
          [
              {
                  "id": "ess_#",
                  "components":
                  [
                    {
                        "component_id": "flexgen_ess_#_hs",
                        "variables":
                        {
                            < in here >
                        }
                    }
                ]
            }
        ]
     }
    }
 ```       
## Update tools design

The update tool allows navigation to any object in a json file.

The navigation is from a designated offset (ArrayIndex) so that complicated layouts can be managed

To nagivate from the root of the json object the ArrayIndex is set to 0.

The "Get" function is then fed a list of nested components. Each subsequent component must be a child of the previous one.
IN this case 

```
ess|asset_instances is an array , we want the array element with "id" = "ess_#".

This array element then has a component called "components" . We want the element with  "component_id" ="flexgen_ess_#_hs".

Having got there we want the "variables" object . 

Once rooted at this spot we can all the required data items.

We can add simple name:value pairs or additional "components" to be nested inside the root componment. 

We can even add Arrays.

 
To help manage Arrays , the array object is discovered and then an item in the array can be selcted by finding a "name=value" field within that array.

Once an ArrayIndex has been discovered, the update tool can add/delete/rename objects from the ArrayIndex designated location.

All these changes are merged into the original text object.


## Proposed Commands

Here are the proposed commands

```
cmd load assets.json
// navigate to ess_# in asset_instances


cmd set idx1 = ArrayIndex from 0 where id = ess_# in [assets,ess,asset_instances] 

// navigate to components
cmd set idx2 = ArrayIndex from idx1 where component_id is flexen_ess_#_id in [components] 

// add a new variable object
cmd AddItem from idx2    "absolute_power_direction_flag" "{}" in [variables]

// add items to the new variable 
cmd AddItem from idx2 "name" "\"Charging when Enabled\""  in    
       [variables,absolute_power_direction_flag]
   
cmd AddItem from idx2 "ui_type" "\"Control\""  in    
       [variables,absolute_power_direction_flag]

cmd AddItem from idx2 "remote_enabled" "true"  in    
       [variables,absolute_power_direction_flag]

// add an options array to the new variable
cmd AddItem from idx2 "options" "[]"  in    
       [variables,absolute_power_direction_flag]

// add an item to the new array
cmd AddItem from idx2 "name" "\"Charge\""  in    
       [variables,absolute_power_direction_flag,options]

//    navigate idx3 to the new item in the options array  where name = Charge 
cmd set idx3 = ArrayIndex from idx2 where name is Charge in 
              "variables","absolute_power_direction_flag","options")

// add item to the option
cmd AddItem from idx3 "value" "true" 

// add an item to the options array
cmd AddItem from idx2 "name" "\"Discharge\""  in    
       [variables,absolute_power_direction_flag,options]

//    navigate idx3 to the name Discharge in the options array  where name = Discharge 
cmd set idx3 = ArrayIndex from idx2 where name is Discharge in 
              "variables","absolute_power_direction_flag","options")

// add item to the option
cmd AddItem from idx3 "value" "true" 

// save it all
cmd save assets.json
```

