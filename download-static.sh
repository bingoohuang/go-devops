#!/usr/bin/env bash

cdn_provider="https://cdnjs.cloudflare.com/ajax/libs"

echo "cdn_provider is $cdn_provider"

mkdir -p static/
cd static

# css resources
rm -fr codemirror.min.css
#curl -LO https://cdn.bootcss.com/codemirror/5.42.2/codemirror.min.css
curl -LO ${cdn_provider}/codemirror/5.42.2/codemirror.min.css

rm -fr jquery-confirm.min.css
#curl -LO https://cdn.bootcss.com/jquery-confirm/3.3.2/jquery-confirm.min.css
curl -LO ${cdn_provider}/jquery-confirm/3.3.3/jquery-confirm.min.css

# javascript resources
rm -fr jquery.min.js
#curl -LO https://cdn.bootcss.com/jquery/3.3.1/jquery.min.js
curl -LO ${cdn_provider}/jquery/3.3.1/jquery.min.js

rm -fr codemirror.min.js
#curl -LO https://cdn.bootcss.com/codemirror/5.42.2/codemirror.min.js
curl -LO ${cdn_provider}/codemirror/5.42.2/codemirror.min.js

rm -fr toml.min.js
#curl -LO https://cdn.bootcss.com/codemirror/5.42.2/mode/toml/toml.min.js
curl -LO ${cdn_provider}/codemirror/5.42.2/mode/toml/toml.min.js

rm -fr jquery.contextMenu.min.js
#curl -LO https://cdn.bootcss.com/jquery-contextmenu/2.7.1/jquery.contextMenu.min.js
curl -LO ${cdn_provider}/jquery-contextmenu/2.7.1/jquery.contextMenu.min.js

rm -fr jquery-confirm.min.js
#curl -LO https://cdn.bootcss.com/jquery-confirm/3.3.2/jquery-confirm.min.js
curl -LO ${cdn_provider}/jquery-confirm/3.3.2/jquery-confirm.min.js
