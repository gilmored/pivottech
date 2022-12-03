# Capstone Project Proposal

## Summary

For my capstone project, I would like to create a CLI tool that takes a user-provided 
VIN (vehicle identification number) and returns information about the vehicle as well as 
opens a browser that takes some of the returned information and opens a browser window that 
takes the user to the page of an auto parts website that displays parts for the requested 
vehicle. This project will integrate with CarMD's API (https://api.carmd.com/member/docs)
and rockauto.com. 

With this project I will:

- Build a REST API
- Interact with an external API to gather information
- Make the returned information easy to digest by the user
- Display good overall coding practices, such as with error handling
- Display good testing practices

## User Stories

### As a user, I would like to enter a VIN and see information as well as search for replacement parts for my vehicle.

**Acceptance Criteria**

Given a VIN is valid, the CLI should provide information regarding the vehicle as well as a direct link 
to search for parts for their vehicle. 

Example: 

    vin 1GNALDEK9FZ108495

would provide: 

    Year: 2015
    Make: CHEVROLET
    Model: EQUINOX
    Engine: L4, 2.4L; DOHC; 16V; DI; FFV
    Trim: LTZ
    Transmission: AUTOMATIC
    https://www.rockauto.com/en/catalog/CHEVROLET,2015,EQUINOX

If the VIN is invalid, a message will appear to inform the user to check their VIN and try again. 
