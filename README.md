# SISCONF 🍱

SISCONF - Sistema de Controle de Ferias (or Grocery Shopping Control System) is designed to help small entrepreneurs manage the orders from their customers. The need arose from
the issues faced by a local entrepreneur which would receive and group the orders from their customers through WhatsApp, creating a spreadsheet in the end to send it to their providers. This lead to inneficient management and
a very time-consuming task of manually checking and sending the orders to each provider known by the entrepreneur. In addition to that, the customers weren't able to track their orders, whether
the order was already placed, on their way or received by the entrepreneur, who acted as middleman throughout the process.

The solution presented by this system is to create a centralized hub for the entrepreneur and the customers, where they can create their orders and track them; and where the entrepreneur can group the orders and automatically create
spreadsheets. The use-case diagram displays all system features

![Casos de Uso](https://github.com/user-attachments/assets/b7e86051-9e1a-40df-a514-dee68dee13e8)

# The Architecture 🏗️

A couple of reasons why we choose to go with microservices architecture, along with other architectural decisions, can be seen down below:

1. It would allow us the scale only the service of creating spreadsheets. Since it would do a lot of I/O operations with multiple rows of a spreadsheet, this is very important. Golang was picked for the language of this microservice because of its performance
2. Our teams were working from different places. This also would allow us to have our own release cycle for each part of the system
3. The decoupled user service (Keycloak) could be used in many different applications if we were the expand the business
4. Choosing a message queue instead of gRPC or REST commmunication is the ideal option here, since the thread could not be block after sending the request to generate a spreadsheet, which could take some time
5. Grafana and Prometheus were utilized for us to monitor our server capacity usage and database metrics, allowing us to make intelligent decisions regarding infrastructure and security in the future

Below, you can find the architectural diagram

![Arquitetura de Software(1)](https://github.com/user-attachments/assets/a7193e6c-d9bc-4456-b39e-f85a1cd492a8)

# The Data Processing Microservice 📊
* Implemented with **Golang** because of its raw performance
* Acted as a consumer, executing jobs created the producer (the SpringBoot application) with round-robin dispatch
* Uploads all the generated spreadsheet to an S3 bucket
* Automatically updates the database utilized by the main API (SpringBoot)

# Entity-Relationship Diagram

![SISCONF - DER](https://github.com/user-attachments/assets/06a14b5e-36ad-486b-b997-ca8bc87f8933)

