## Web scrapping using producers y consumers


**Description:** 
It is required to store the relationship between films and actresses in a database.
The films in which each actress has worked are defined on each web page of the filmography in which the actress has acted.

- Diagram:


![prodcondiagram](https://github.com/Pjvl99/Con-Prod/assets/61527863/3174c81b-0a2e-4146-b938-ad53ccd5e5ba)

### Tools:
- Go: To create and manage consumers and producers
- Docker
    - mysql: Store all data

## Project Execution
#### Database - Docker Mysql
```
docker-compose up
``` 

#### GO - Producers and consumers
```
./pcb #Producers #Consumers
```
In #Producers we put the number of producers and in #Consumers the number of consumers. To run this command docker must be running.

### Execution
The idea is to create a number of producers who collect the data for each actress and pass it to consumers without a race condition occurring that causes the data to be repeated.

![img1](https://user-images.githubusercontent.com/61527863/169353015-8e18b421-fc4e-4f4c-aeff-4f6e19153791.jpeg)

![img4](https://user-images.githubusercontent.com/61527863/169354683-68796bdd-935c-4026-af43-6893058016b9.jpeg)

![img5](https://user-images.githubusercontent.com/61527863/169354693-b7e2aa4a-4b3b-4ca9-98ae-b7e345f41d8c.jpeg)

![img6](https://user-images.githubusercontent.com/61527863/169354743-bef545b1-45b0-4497-9e5d-04aa796a8cac.jpeg)

In the database you can see the results, without any repeated actress data.

![img2](https://user-images.githubusercontent.com/61527863/169353632-0404ee0a-17f4-4e96-8b68-fe7c730d2426.jpeg)

![img3](https://user-images.githubusercontent.com/61527863/169354305-68f2977e-48ef-437b-80a4-234268f83136.jpeg)

