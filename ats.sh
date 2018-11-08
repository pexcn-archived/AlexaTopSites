#!/bin/bash -e

export ATS_ACCESS_KEY="YOUR_ACCESS_KEY"
export ATS_SECRET_KEY="YOUR_SECRET_KEY"

/root/ats/ats $ATS_ACCESS_KEY $ATS_SECRET_KEY CN 1 1000 > 1k.xml
/root/ats/ats $ATS_ACCESS_KEY $ATS_SECRET_KEY CN 1001 1000 > 2k.xml
/root/ats/ats $ATS_ACCESS_KEY $ATS_SECRET_KEY CN 2001 1000 > 3k.xml
/root/ats/ats $ATS_ACCESS_KEY $ATS_SECRET_KEY CN 3001 1000 > 4k.xml
/root/ats/ats $ATS_ACCESS_KEY $ATS_SECRET_KEY CN 4001 1000 > 5k.xml

cat 1k.xml 2k.xml 3k.xml 4k.xml 5k.xml | sed -n 's|.*<aws:DataUrl>\(.*\)</aws:DataUrl>.*|\1|p' > top-5000.txt



git clone https://github.com/pexcn/AlexaTopSites.git -b gh-pages --depth 10 gh-pages
pushd gh-pages

find . ! -path . | grep -v ".git" | xargs rm
mv ../*.xml ../*.txt .

git add --all
git commit -m "`date +'%Y-%m-%d %T'`"
git push "https://USERNAME:PASSWORD@github.com/OWNER/REPO.git" gh-pages
popd
rm -r gh-pages
