require('dotenv').config({ path: '.env' });
require('child_process').spawn('yarn', ['dev'], { stdio: 'inherit', shell: true }); 