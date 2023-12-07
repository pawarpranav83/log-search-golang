<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a name="readme-top"></a>
<!--
*** Thanks for checking out the Best-README-Template. If you have a suggestion
*** that would make this better, please fork the repo and create a pull request
*** or simply open an issue with the tag "enhancement".
*** Don't forget to give the project a star!
*** Thanks again! Now go create something AMAZING! :D
-->



<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->


<!-- PROJECT LOGO -->
<br />
<div align="center">

<h3 align="center">Log-App-Golang</h3>

  <p align="center">
    Log Ingestor and Log Search App
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#resources">Resources</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

Log ingestor system that can efficiently handle vast volumes of log data, and offer a simple interface for querying this data using full-text search or specific field filters.

<p align="right">(<a href="#readme-top">back to top</a>)</p>


### Built With

* Golang
* Gin
* MongoDB (Atlas)
* HTML (For UI)
* Docker

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

This is an example of how you may give instructions on setting up your project locally.
To get a local copy up and running follow these simple example steps.

### Prerequisites

Need to have either docker or go installed.

### Installation

#### Docker Based


1. Clone the repo
   ```sh
   git clone https://github.com/dyte-submissions/november-2023-hiring-pawarpranav83.git
   ```

2. Run docker-compose yaml file. It will build the image and run it on port :3000
   ```sh
   docker compose up
   ```


#### Golang Based


1. Clone the repo
   ```sh
   git clone https://github.com/dyte-submissions/november-2023-hiring-pawarpranav83.git
   ```

2. Run
   ```sh
   go build .
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage


### Log Ingestor

Perform a POST request with Log data on http://localhost:3000/

**Note** - The format should be as per the given sample, if any of the fields are nil, the server will give an error.
This validation is performed using JSON and binding tags (provided by gin).

It would return an InsertedID, for successful insertion in the database.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Query Interface

HTML UI can be accessed with GET request on http://localhost:3000/

Uses an HTML Form with fields for searching, and on submit, makes a POST request to http://localhost:3000/search, which displays the docs obtained.
One field for Full-Text Based Search - Text-Search.

Other fields act as filters.



**Created an Aggregation Pipeline** - 
1. First Stage - Full-Text Based Search on all documents of a collection (created a search index using atlas search).

2. Second Stage - Matching / Filtering the documents obtained from the first stage based on the input filter parameters. Equal to operation, but for message we use the text-based search (created a text index for message), but can only be used when Text-Search field is empty (since MongoDB only allows text-based search on first stage).

3. Third Stage - Limit the number of documents to five.



<!-- ROADMAP -->
## Requirements

### Log Ingestor

- [x] Develop a mechanism to ingest logs in the provided format. Performed Validation using struct tags and the ShouldBindWithJSON method provided by Gin.
- [x] Ensure scalability to handle high volumes of logs efficiently. Gin processes requests concurrently, the server can handle high volumes of logs.
- [x] Mitigate potential bottlenecks such as I/O operations, database write speeds, etc. Since only text-based search indexes are added, no bottlenecks in write speeds, and since requests process concurrently so no blocking because of I/O operations.
- [x] Make sure that the logs are ingested via an HTTP server, which runs on port 3000 by default.

### Query Interface
 - [x] Offer a user interface (Web UI or CLI) for full-text search across logs. Implemented HTML Form for searching logs.
 - [x] Filters. Have added filters for all fields and can be used by specifying a value for a field in their corresponding input field.
 - [x] Aim for efficient and quick search results. Have search-based indexes that result in a quick text-based search.

### Advance Features
- [x] Implement search within specific date ranges. Using gte and lte operators in match stage for timestamp. Also timestamp is stored as value using time.Parse method in golang.
- [x] Allow combining multiple filters. Multiple filters can be combined, if not null those fields are added in the match stage filters. Combined using AND operator.
- [ ] Provide real-time log ingestion and searching capabilities. Can be implemented using websockets and refreshing the result whenever the server emits an event. We just upgrade to web socket protocol from HTTP (can use gorilla/websockets) and then, specify events that server and client listens to.
- [ ] Implement role-based access to the query interface. Can be implemented by creating an authentication process with role assigned by an admin, and adding a middleware function before the search route to check whether that specific role can specify certain filters or not.

I know it's easier said than done, but my major examinations start on Monday, that's why I wasn't able to implement the last two features.

However, I have implemented websockets and authentication with role-based restrictions in Golang in my other project Golang-Chat-App (https://github.com/pawarpranav83/golang-chat).

Also have created a **Dockerfile** so that the image can be used in various cloud-based solutions for scaling up.

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- CONTACT -->
## Contact

Pranav Pawar - [@pawarpranav83](https://www.linkedin.com/in/pranav-pawar-b54954242/) - pawar.pranav83@gmail.com

Phone Number - 8104273543

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ACKNOWLEDGMENTS -->
## Resources

* MongoDB Go Driver Docs - https://www.mongodb.com/docs/drivers/go/current/
* Gin Docs - https://github.com/gin-gonic/gin/blob/v1.9.1/docs/doc.md

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: [https://linkedin.com/in](https://www.linkedin.com/in/pranav-pawar-b54954242/)