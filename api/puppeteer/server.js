const express = require('express');
const app = express();
const port = 3939;
const { exec } = require('child_process');

const service = '{"ID":"777","Name":"maxtraffic","LogoURL":"nourl","Patterns":[{"ID":"scriptpatternid","Type":"resource","Value":"*mt.js*","Mandatory":true}]}';

app.get('/', (req, res) => {
    exec(`node analyze.js https://maxtraffic.com '${service}'`, (err, stdout, stderr) => {

        console.log(err, stdout, stderr)
        res.send(`Hello World! ${err}, ${stdout}, ${stderr}`)
    });
});

app.listen(port, () => console.log(`Example app listening on port ${port}!`));
