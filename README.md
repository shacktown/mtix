## HyperLedger Fabric BlockChain Movie Sample application
 
 It's Movie Night on the Blockchain!

 This "Movie Night" HyperLedger Fabric (Fabric) sample application provides a blockchain implementation for a Movie Cinema company to manage schedules, sales and inventory for one or more theater locations. This includes scheduling movie showings, stocking concessions, selling tickets, selling refreshments and more. 

 Capabilities provided by BlockChain and HyperLedger fabric include: 
 - Secure, authenticated and authorized access
 - Sales and inventory tracking and management
 - Up to the minute dynamic data for trend analysis, forecasting and machine learning
 - The potential for provenance and validation through the entire sales chain including movie licensing, supply tracking, etc. 

 The sample application could be easily extended to support a consortium of cinema industry vendors and suppliers. This would provide an open, secure system for managing the end to end supply chain.

 ----------
 ### The Use Case Scenario

 The application models a Movie Cinema company with one or more locations. Each theater location provides one or more auditoriums with a specified seating capacity and other features (e.g. 3D, Surround Sound, etc.). Each location can also provide concessions such as soda, water and popcorn. The application manages  inventory, show schedules and sales for each location.

----------
 ### The Transactions

 ***Theater Managers*** schedule movie showings and update concession inventories. Blockchain transactions for these operations include:
- ScheduleMovieShowing
- StockConcession

***Theater Staff*** use provided client applications to sell movie tickets and concessions. Blockchain transactions for these track available seats for each movie showing as well as concession inventories, and include:
- BuyTix
- BuyConcession
- TicketsAvailable
- ConcessionsAvailable
- ExchangeConcession
----------

### The Data Model
The "Movie Night" sample application uses a simplified data model to capture key informational attributes of each modeled entity. The data model could be easily extended to capture and track additional detail. 

<u>Theaters</u>

| ID (key) | Name | Address | Phone | Web URL | #Auditoriums |
| :---: | :---: | :---: | :---: | :---: | :---: | 
| MNcinema | Movie Night Cinema1 | 123 1st Ave, Apex, NC 27502 | 555 555-5555 | www.movienight1.com | 8 |
| NTcinema | Nighttime Cinemas  | 456 2nd Ave, Austin, TX 12345 | 444 444-4444 | www.nighttimecinema.com | 12 |

<u>Movies</u>

| ID (key) | Title | Length | Rating | ReleaseDate | 
| :---: | :---: | :---: | :---: | :---: | 
| BttF | Back to the Future | 120 | PG-13 | 05-05-1984 | 
| RotLA | Raiders of the Lost Ark | 115 | PG-13 | 05-05-1981 | 

<u>Scheduled Showings</u>

| ID (key) | Theater ID | Movie ID | Auditorium ID | Show Time | Price | Capacity | Tickets Sold | 
| :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | 
| show7NA20 | MNcinema | BttF | Hall7 | 2019-02-29T21:00:00-05:00 | 11.00 | 120 | 47 |
| show4TA44 | NTcinema | RotLA | Hall11 | 2019-02-31T19:00:00-06:00 | 14.00 | 90 | 57 |

<u>Ticket Sales</u>

| Txn ID | Show ID | Quantity | Purchase Time | Gross Revenue | Seller ID | 
| :---: | :---: | :---: | :---: | :---: | :---: |
| N7423A | show7NA20 | 4 | 2019-02-05T02:00:00-05:00 | 44.00 |  MNcinema-Window2 |
| X5JC72 | show4TA44 | 2 | 2019-02-14T14:00:00-06:00 | 28.00 |  NTcinema-Window4 | 

<u>Concessions</u>

| Theater ID | Concession ID | Inventory (#units) | Price |
| :---: | :---: | :---: | :---: | 
| MNcinema | Popcorn | 257 | 8.00 |
| MNcinema | 16oz-Water-Bottle  | 421 | 5.00 |

<u>Concession Sales</u>

| Txn ID | Theater ID | Concession ID | Quantity | Purchase Time | Gross Revenue | 
| :---: | :---: | :---: | :---: | :---: | :---: |
| C5724Z | NTcinema | 16oz-Water-Bottle | 2 | 2019-02-05T02:00:00-06:00 | 10.00 |
| C5725A | NTcinema | Popcorn | 1 | 2019-02-14T04:00:00-06:00 | 8.00 |

----------

### Coming Soon
- Instructions for installing and running the application
- Descriptions for each source file
- Suggestions for contributions and enhancements
