function elem(tag, clazz, html) {
    // TODO: exists something standard for this?
    var element = document.createElement(tag);
    if(clazz)
        element.className = clazz;
    if(html)
        element.innerHTML = html;
    return element;
}

function renderHeader() {
    var hh = document.getElementsByTagName("header");
    if(!hh.length) {
        console.error("Couldn't find header element");
        return;
    }
    var h = hh[0];
    while(h.firstChild) {
        h.removeChild(h.firstChild);
    }

    var hadd = function(code){
        h.insertAdjacentHTML('beforeend', code);
    }
    hadd('<a href="/"><img src="/default_20200205_/img/wheel_48x48.png" width="48" height="48" class="header_picto" /></a>');
    hadd('<h1><a href="/">Programming-Idioms</a></h1>');
}

function renderFooter() {
    var footerz = document.getElementsByTagName("footer");
    var footer = footerz[0];
    footer.insertAdjacentHTML('beforeend', `
    <div>
		All content <a href="http://en.wikipedia.org/wiki/Wikipedia:Text_of_Creative_Commons_Attribution-ShareAlike_3.0_Unported_License" rel="license">CC-BY-SA</a>
    </div>
    <div>
		<a href="/about#about-block-language-coverage" title="Coverage grid"><img src="/default_20200205_/img/coverage_icon_indexed.png" class="coverage" /></a>
	</div>
    <div>
		<a href="/about" class="about-link">?</a>
	</div>`);
}

function highlightDie(){
    var die = document.querySelector(".die");
    var src = die.src;
    var hsrc = src.replace("die.svg", "die_highlight.svg");
    die.onmouseover=function(){this.src=hsrc;};
    die.onmouseout=function(){this.src=src;};
}

//
// Execution!
//

renderHeader();
renderFooter();
highlightDie();