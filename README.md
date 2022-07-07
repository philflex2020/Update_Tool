### Update Tool

p wilshire  07/01/2022

## Progress

Investigating a better way to do this.
jparse_test2.go  is a json parser with a difference.

It allows you to read in a text file and apply json parsing techniques to it without changing the source file.
It works on a big byte array and extracts pointers to data  using json references.
We can then extract , update etc as needed . we end up with a new byte string that we can save as a newfile.

This means we can have a set of instructions and apply those to the file and maintain the original file string format.


....


p wilshire  06/22/2022


The update tool is tasked with appying a complete release update to a system.
Currently the update process can take several days, or weeks,  to complete. 
This is due to the complexity of the system configuration combined with the customer needs to get the upgrade completed within a short timeslot.
Config items can easily be missed or not performed correctly.

The Update tool is designed to address the update operation and streamline the whole process.
Note that the update does not cover just the system configs but all the OS tool configs for databases and system start up.

Any deficiencies in the Update Tool will be fixed as the tool is developed giving an improving solution to the Update problem.


## What's needed

The FlexGen release process provides a working set of systems using development configurations.
These configurations often have to be changed as the development system is deployed to a particular site.
Some of these changes have to be revised if the deployment is to be repeated on a similar but different site.
The upgrade problem lies in repeating all the customisation changes as a new "development" release is applied to each site.

The upate tool needs a way to capture the deployment modifications required for a particular site and then allowing those changes to be adapted 
to make the system work on a different site.

Consider the two TX100 sites. 
The 10.2 Software release will have to be tested against a NorthFork site configuration. 
In the past this has been done partially in the lab (or Gauntlet).
The development changes are modified to create a working system in the simulation or development environment. 
The completed configurations are then applied to the Northfork site where a number of additional adjustments may need to be made. 
This is due to differences between the Development simulation and the actual site.
The update_tool has to capture these changes. 

When the same system is deployed to the BatCave site we have two groups of configuration changes to be applied.
One group pertains to the Simulation -> Site deployment changes 
The other group is related to the Site->Site differences.


The update_tool is intended to capture both sets of changes to ensure a smooth transition of the deployment to the  different site.

## Config Changes Development -> Release

The lab test environment is intended to test the appliction operation against a Twins model of the system in a development environment.
That twins model can, in theory, be adjusted to reflect the characteristics of a target installation.
This functionality needs to be enhanced using metric, Fims_echo and / or Ess_sim tools.

The Bms, Pcs operations can be customised to match that of a particular site combination.
This covers the development to the deployment environment. 
It also covers changes to transition from development mode to deployment mode.
It will work against the simulation of a site as far as that simulation can represent that site's operation.

These changes are covered pretty well by the current design flow. The configuations for 10.2 Northfork can be captured during the release cycle.

## Config Changes Release -> site Problem #1

These changes have to be made to the system to compensate for the deficiencies in the simulation environment compared to the actual site.
Adjustments are made to  transition from the Development release to the site deployment.
This is a time consuming operation, visible to the customer. It can take time to identify the correct config changes to make the system work properly.
These changes cannot,currently,  be applied to the release configurations and tested in a simulation environment.

## Config Changes site -> site Problem #2

Even if all the config changes are captured from the deployment on Site #1, these changes cannot simply be applied to Site #2.
There is a mixture of Release->Site changes and Site Customization changes.
The examples of the nature of this problem  is the database running out of space because of incorrect retention policies from the other site are left in the storage config.

## Manage two kinds of changes

The update_tool has to be made aware of the two types of config changes. Simulation_to_site and Site_to_site.

## Auto_config tool

This is a system designed to partially solve this problem.
The config requirements for development testing for a given release on a specifed site are clearly identified with this tool.
THe split between design config and site specific config are also be clearly identified in this system.

The update_tool has to specifically capture those config requirements but also easily identify the site to site differences.

### Changes due to release updates.

This category of config changes is needed to allow a working site configuration to be modified to accomodate design changes in release cycles.
If a particular operation requires a different set of operating controls or parameters the config for the old operation needs to be replaced by 
the config required for the new operation.

However, as these config changes are applied, the system still needs to accomodate the Simulatio -> site changes and the site->site changes.

### How an update is applied.

## Development

We have a running configuration for the 10.1 release registered with the Update_Tool.
This means that all the configs monitored by the Update_Tool have been saved 

A new 10.2 release is then commissioned in the lab.

The Update_Tool will identify all changes in the lab configuration.

For example 

Site Controller Metrics config change detected, please classify.

Each change will then be classified as in the site deployment (described below).

Each change can be accompanied by a test process to ensure the correct operation of the change.

During system testing, the config changes can be applied in layers.

The Update_Tool can then capture each layer and assemble a series of patches required to complete the update.
So changes for CATL BMS systems can be isolated from changes for Risen BMS systems.

After MVP we may need to fine tune this progressive patch approach.

The mongodb database will be used to capture the config changes.
This ensures that the same changes can be managed on site.


## Site deployment

This should be mainly covered by the Lab Simulations.
If any differences between the Lab Simulation and the system in the field. The update_tool should provide a change notice.
This should cover, scaling differences, variable name changes, default values etc.
The change notice should track any changes between the update distribution and the applied update in the field.
Consider this as a sort of patch.
   "Whenever we apply 10.3 to CATL and Power Electronics we need to make these changes"

This will subivide down to the following types of changes.
   This patch will apply to all instances where CATL BMS systems are used 
   This patch will apply to all instances where PE PCS systems are used 
   This patch will apply to all instances where CATL BMS units are used together with PE PCS systems. 


Also this
"Whenever we apply 10.3 to CATL and Power Electronics on Tx100 configurations we need to make these changes"
   This patch will apply to all instances where CATL BMS units are used together with PE PCS systems on TX100 design systems. 

And this
"Whenever we need to apply the system on Northfork to BatCave we need these changes"
    This patch is applied to make the system work at NorthFork
    This patch is applied to make the system work at BatCave


So the deployment engineer will adjust the configuration for a particular site and then get the Update_Tool to capture the changes.
The change will then have to be classified and recorded.

Using a mongodb containing patches  on an external system ( or container) allows the whole site config to be managed by the Update_Tool.


## More application details.

When a change is to be applied to a system the following steps are taken.

 - check that the unpatched system conforms to a state expected by the patch application. We'll have to take out, or wildcard , any site specific information.
 - apply the patch.
 - update any site specific adjustments.
 - perform any tests to ensure that the operation is as expected after the patch has been applied.

All results from this process are recorded in detail. These details provide an update record.
