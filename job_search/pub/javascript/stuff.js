(function() {
    window.onload = function() {

        var form = document.querySelector("#search-form")

        form.addEventListener("submit", function(e) {
            e.preventDefault()

            document.getElementById("results").innerHTML=""

            var x = new XMLHttpRequest()
            x.onreadystatechange = function() {
                if(x.readyState === 4 && x.status === 200) {
                    obj = JSON.parse(x.response)
                    makeHTML(obj)
                    addStatusUpdate()
                    changeColor()
                    getResInfo()
                }
            }
            x.open("POST", "/post")
            x.send(new FormData(form))

        });

        var filter = document.querySelector("#filter-kw")
        filter.addEventListener("keyup", function() {
            filter = document.querySelector("#filter-kw")
            var companies = document.querySelectorAll(".company")
            var descriptions = document.querySelectorAll(".desc-text")

            filterwords = filter.value.split(",")
            for(var i = 0; i < descriptions.length; i++) {
                for(var j = 0; j < filterwords.length; j++) {
                    if(companies[i].innerText.toUpperCase().includes(filterwords[j].toUpperCase())){
                        companies[i].closest(".job-box").style.display = null
                        break
                    } else if(descriptions[i].innerText.toUpperCase().includes(filterwords[j].toUpperCase())) {
                        descriptions[i].closest(".job-box").style.display = null
                        break
                    } else {
                        descriptions[i].closest(".job-box").style.display = "none"
                    }
                }
            }
        });

        var hide = document.getElementById("hide-no-interest")
        hide.addEventListener("click", function(){
            hideNotInterested(this)
        });
    }
})()

function hideNotInterested(hide) {
    if(hide.type == "checkbox") {
        if(hide.checked == true){
            var appliedstatus = document.getElementsByClassName("status")
            for(var i=0; i < appliedstatus.length; i++){
                if(appliedstatus[i].value == 3) {
                    appliedstatus[i].closest(".job-box").classList.add("checkbox-hide")
                }
            }
        } else {
            var hidden = document.getElementsByClassName("job-box")
            for(var i=0; i < hidden.length; i++) {
                hidden[i].classList.remove("checkbox-hide")
            }
        }
    } else if(hide.type == "range"){
        var hidecheck = document.getElementById("hide-no-interest")
        if(hide.value == 3 && hidecheck.checked == true){
            hide.closest(".job-box").classList.add("checkbox-hide")
        }
    }

}



// Appends each job to the page in this style
// <div id="j" + index> class="job-box">
//     <div class="info-box">
//         <h3 class="title"></h3>
//         <h3 class="company"></h3>
//         <h3 class="date"></h5>
//         <input type="range" class="status">
//     </div>
//     <div class="description">
//         <div></div>
//         <p class="apply"></p>
//     </div>
// </div>
function makeHTML(obj){
    for(var i=0;i<Object.keys(obj).length;i++){
        var results = document.querySelector("#results")
        
        // <h1 class="title"></p>
        var title = document.createElement("h3")
        title.setAttribute("class", "title")
        var title_text = document.createTextNode(obj[i].jobTitle)
        title.appendChild(title_text)

        // <p class="company"></p>
        var company = document.createElement("h3")
        company.setAttribute("class", "company")
        var company_text = document.createTextNode(obj[i].companyName)
        company.appendChild(company_text)

        // <p class="date"></p>
        var date = document.createElement("h3")
        date.setAttribute("class", "date")
        var date_text = document.createTextNode(obj[i].dateAdded)
        date.appendChild(date_text)

        // <div class="info"></div>
        var job_info = document.createElement("div")
        job_info.setAttribute("class", "info-box")

        var status = document.createElement("input")
        status.setAttribute("class", "status")
        status.setAttribute("type", "range")
        status.setAttribute("min", "1")
        status.setAttribute("max", "3")
        status.setAttribute("value", obj[i].appliedStatus)
        var statspan = document.createElement("span")
        statspan.setAttribute("class", "stat-span")
        statspan.appendChild(status)

        // <div class="info">
        //     <h3 class="company"></h3>
        //     <h3 class="title"></h3>
        //     <h3 class="date"></h3>
        // </div>
        job_info.appendChild(company)
        job_info.appendChild(title)
        job_info.appendChild(date)
        job_info.appendChild(statspan)
        
        // <div class="description"></div>
        var description = document.createElement("div")
        description.setAttribute("class", "description")
        var desc = obj[i].jobDescription
        description.innerHTML = "<div class='desc-text'>" + desc + "</div>"

        // <p class="apply">link</p>
        var apply = document.createElement("a")
        apply.setAttribute("class", "apply")
        apply.setAttribute("href", obj[i].applyLink)
        var apply_text = document.createTextNode("Apply")
        apply.appendChild(apply_text)
        description.appendChild(apply)



        // <div id="j#"></div>
        var jbox = document.createElement("div")
        jbox.setAttribute("id", obj[i].key)
        jbox.setAttribute("class", "job-box")

        jbox.appendChild(job_info)
        jbox.appendChild(description)

        results.appendChild(jbox)
    }
    // Add open and closing of description when you click on the info box
    // and scrolls the page so the chosen job starts at the top
    var jobs = document.getElementsByClassName("info-box")
    for (var i=0; i<jobs.length; i++){
        jobs[i].addEventListener("click", function(e) {
            this.nextSibling.style.display = this.nextSibling.style.display == "" ? "block" : ""
            window.scroll({top: window.pageYOffset + this.parentNode.getBoundingClientRect().top, left: 0, behavior: "smooth"})
        });
        
    }
}

function addStatusUpdate(){
    var appliedstatus = document.getElementsByClassName("status")
    var appliedstatusarea = document.getElementsByClassName("stat-span")
    for(var i=0; i < appliedstatus.length; i++){
        appliedstatusarea[i].addEventListener("click", function(e){
            event.stopPropagation();
        });

        appliedstatus[i].addEventListener("mouseup", function(e){
            e.stopPropagation();
            var x = new XMLHttpRequest()
            x.onreadystatechange = function() {
                if(x.readyState === 4 && x.status === 200) {
                }
            }
            x.open("POST", "/updateStatus")

            var formData = new FormData()
            formData.append("status", this.value)
            formData.append("key", this.closest(".job-box").id)

            x.send(formData)
            changeOneColor(this, this.closest(".job-box").id)
            getResInfo()
            hideNotInterested(this)
        });
    }
}


function changeOneColor(as, i) {
    var infobox = document.getElementsByClassName("info-box")

    if(as.value == 1){
        infobox[i-1].style.background = "#a1c22d"
    } else if(as.value == 2) {
        infobox[i-1].style.background = "#9eb9c3"
    } else if(as.value == 3) {
        infobox[i-1].style.background = "#c02942"
    }
}

function changeColor() {
    var infobox = document.getElementsByClassName("info-box")
    var statusval = document.getElementsByClassName("status")

    for(var i=0; i<statusval.length; i++){
        if(statusval[i].value == 1){
            infobox[i].style.background = "#a1c22d"
        } else if(statusval[i].value == 2) {
            infobox[i].style.background = "#9eb9c3"
        } else if(statusval[i].value == 3) {
            infobox[i].style.background = "#c02942"
        }
    }
}

function getResInfo() {
    var statusval = document.getElementsByClassName("status")
    var app = 0
    var nointerest = 0
    var other = 0

    for(var i=0; i<statusval.length; i++){
        if(statusval[i].value == 1){
            app++
        } else if(statusval[i].value == 2) {
            other++
        } else if(statusval[i].value == 3) {
            nointerest++
        }
    }

    appliedcount = document.getElementById("res-info-1")
    othercount = document.getElementById("res-info-2")
    nointerestcount = document.getElementById("res-info-3")

    appliedcount.innerText = app
    othercount.innerText = other 
    nointerestcount.innerText = nointerest
}