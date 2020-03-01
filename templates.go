// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bounce
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

const htmlTemplate string = `
<!doctype html>

<html lang="en">
<head>
  <meta charset="utf-8">

  <title>bounce - Small utility to assist with building HTTP endpoints</title>
  <meta name="description" content="bounce - Small utility to assist with building HTTP endpoints">
  <meta name="author" content="atc0005">

  <!--
	  https://www.w3schools.com/css/css_table.asp
	  http://web.simmons.edu/~grabiner/comm244/weekfour/code-test.html
  -->
  <style>

  body {
	  padding: 0em 0.5em 0em 0.5em;
  }

  table {
	border-collapse: collapse;
	width: 100%;
  }

  th, td {
	text-align: left;
	padding: 8px;
  }

  tr:nth-child(even){background-color: #f2f2f2}

  th {
	background-color: #4CAF50;
	color: white;
  }

  code {
	background-color: #eee;
	border: 1px solid #999;
	display: block;
	padding: 0.5em;
  }

  </style>

</head>
<body>

<h1>Welcome!</h1>

<p>
  Welcome to the landing page for the bounce web application. This application
  can be found at <a href="https://github.com/atc0005/bounce">
  github.com/atc0005/bounce</a>. There you can check for newer versions or
  submit bug reports. Feedback is welcome.
</p>

<h2>Purpose</h2>

 <p>
  This application is primarily intended to be used as a HTTP endpoint for
  testing webhook payloads. Over time, it may grow other related features
  to aid in testing other tools that submit data via HTTP requests.
</p>

<h2>Supported endpoints</h2>

<p>
  The list of links below are the currently supported endpoints for this
  application:
</p>

<table>
  <tr>
    <th>Name</th>
    <th>Pattern</th>
    <th>Description</th>
    <th>Allowed Methods</th>
  </tr>
{{range .}}
  <tr>
    <td><code>{{ .Name }}</code></td>
    <td><code>{{ .Pattern }}</code></td>
	<td><code>{{ .Description }}</td>
	<td>{{range .AllowedMethods}}<code>{{ . }}</code> {{end}}</td>
  </tr>
{{else}}
<tr>
  <td><code>Failed to parse routes</code></td>
  <td><code>N/A</code></td>
  <td><code>N/A</code></td>
  <td><code>N/A</code></td>
</tr>
{{end}}
</table>

<h2>Feedback</h2>

<p>
  NOTE: The primary developer of this application is new to Go, so there are
  likely many rough edges. Please let us know what problems you encounter
  so that we may work to resolve them in a future release.
</p>



</body>
</html>
`
