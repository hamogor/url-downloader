# Backend Engineering Exercise

In this exercise you will develop a linux daemon with the following features:

- It should expose a minimal HTTP server. You don't need to write a full HTTP Server on your own; rather, you should use the Golang std lib net/http package or a framework library of your choice (like Go Gin which is what we use internally).
- The web server should serve the required REST APIs.
- The API should expose a method able to receive a request to add a website URL to an internal list of objects (example URL: http://www.example.com/).
- The API should expose an endpoint which is able to retrieve the latest 50 URLs sent through the previous method (sorted from newest to oldest or from the smallest to the biggest, upon user request) and a counter that would show how many times that specific URL have been submitted to the API since the program started.
- Upon submission, each URL should be downloaded (GET request) but no more than 3 downloads at the same time should be executed. If the initial download of the page fails, throw the URL away. If the download is successful, the URL should be stored and used later.
- Create a background process which is executed every 60 seconds. This process must get the 10 most submitted/requested URLs from the ones that have been submitted and try to fetch the URL again, measuring the time it took to download it. All the download operations should happen in parallel with a concurrency factor of 3 - so no more than 3 GET requests should happen at the same time.
- Collect all the download times, successful downloads counter and failed downloads counter and log them all on the stdout when the previous batch process completes.
- Ideally, the API should be ready to be used by a single page JS application running in another product (implementing a UI is out of scope).

# Tools, Libraries and Frameworks

Your solution should build to single runnable executable and be provided with an external yaml configuration file, containing parameters which you think should be exposed. Your application will be reviewed on an up-to-date version of Golang (1.23) on linux. Ideally it will be Dockerised but this is not a strict requirement. There is no need to deploy or create any infrastructure for this task.

Other than Golang, the choice of libraries, frameworks or tools used to develop the application is left open to you. We encourage you to use options you’re familiar with, and make choices that you’d feel comfortable explaining and justifying in the next interview call. Provide evidence of testing your code, but we do not expect 100% coverage.

Please ensure the code you write can be shared with our team through a public repository (e.g. GitHub). We'll review your code, and then ask you to a Zoom interview so we can talk through your ideas.

We value your time and appreciate the effort you put into this task. While there's no expectation to produce a fully polished solution, we encourage you to focus on demonstrating your approach and thought process. If you feel your unfinished code sufficiently showcases your skills and understanding, please include notes in your README.md explaining what you would improve or do differently with more time. We'd be happy to review it.

If you have any problems or issues working with the challenge, please contact us.

# Assessment

We will assess your application based on the following criteria:

* How clean, modular and maintainable the code is.
* Suitability of tools, libraries and frameworks used (for both the app itself and any build processes involved).
* Testing of the code.
* Documentation.