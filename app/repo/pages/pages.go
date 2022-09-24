package pages

const DefaultPageContent = `
	<h1>URL shortener</h1>

	<form action="/" method="post">
	  <h3><label for="lurl">Long URL:</label></h3>
	  
	  <h3> <input type="text" name="lurl" value="%v"/> </h3>
	  <br>
	  <h3> <label>Short URL:</label> </h3>
	  <br>
	  <label>%v </label>
	  <br>
	  <h3> <label>Statistics URL:</label> </h3>
	  <br>
	  <label>%v</label>
	  <br><br>
	  <input type="submit" value="Generate short URL and Statisctics URL">
	</form>
	`
const ShortLinkUrl = `http://localhost:8000/sl/`
const StatUrl = `http://localhost:8000/stat/`
