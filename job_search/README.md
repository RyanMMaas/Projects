# One Stop Job Search

# Why
This project was a way for me to show off using various different languages on a larger project. It was also a project that I enjoyed doing since I am able to use it frequently and I can constantly find new ways to improve on it.

# Keywords
>Go, golang, JavaScript, HTML, CSS, web scraping, JSON, AJAX 

# Features
* [x] Searching [Glassdoor](https://www.glassdoor.com/)
* [x] Searching [Indeed](https://www.indeed.com)
* [x] Filtering results
* [x] Hiding results marked 'not interested'


# Usage
```bash
go run main.go
```
>Visit 127.0.0.1:9000 in browser to visit the site

>ctrl+c to end

# Issues
* Indeed company name shows up wrong

# After Thoughts
This project was fun and useful to me as well as a great learning opportunity. Most of the troubles I had stemmed from trying to find out how the different websites worked so that I could gather the data from each page automatically.

After more research into web scraping I looked into the *robots.txt* of each site and noticed that some of the pages I was scraping were disallowed so I will no longer be using/updating this project, unless I use the API's of the sites.

If I were to do this project again I would rather use the API's of each site as it would be much easier than the method I have currently employed, although I did enjoy doing it this way as a learning experience. Another change I would make would be storing everything in a database rather than a JSON file and making the page load a limited number of postings at a time, as I believe both of these would speed up the site.