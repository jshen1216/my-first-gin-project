# This is my first gin project
This project will include following functions:
- get raw data in json from local mySQL database named db_log
- get raw data in json from local mySQL database named db_log by searching in time period, IP, and keyword
- save the searching result in csv file in folder named tempcsv and the file will be named by account name
- api for downloading the latest search result by csv file

New functions added on 28 April 2022:
- get raw data in json from Elasticsearch in VM which index named accesslog
- get raw data in json from Elasticsearch in VM which index named accesslog by searching in time period, IP and keyword
- will also saved temp csv file to download when searching in ES as well

New functions added on 1 May 2022 (dev branch):
- use swaggo to add API documents
- swagger ui is set up 
-- APIs on swagger ui: get all data from mySQL & ES
-- get data under conditaional serach from mySQL & ES
-- download csv file
