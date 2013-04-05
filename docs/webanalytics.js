jQuery(document).ready(function() {
  function getPagePosition(element){
    var rect = element.getBoundingClientRect();
    return {
      left: rect.left + document.body.scrollLeft,
      top: rect.top + document.body.scrollTop,
      width: rect.width,
      height: rect.height
    };
  }

  var screenHeight = screen.height;
  var screenWidth = screen.width;
  var pageView = {
    "domain" : document.domain,
    "url" : document.URL,
    "userAgent" : navigator.userAgent,
    "screenHeight" : screen.height,
    "screenWidth" : screen.width
  }
  jQuery.ajax({
    type: "POST",
    url: "http://192.168.122.150:8080/page-views/",
    data: JSON.stringify(pageView),
    dataType: "json",
    contentType: "application/json; charset=utf-8"
  });

  jQuery('a').click(function(event) {
    event.preventDefault();
    href = jQuery(this).attr('href');
    // TODO investigate why I added docHeight and docWidth. I think I"ll have to do math.
    //Document height is used with the position of dom elements.
    var docHeight = jQuery(document).height();
    var docWidth = jQuery(document).width();

    rect = this.getBoundingClientRect();
    pagePosition = getPagePosition(this);
    var hrefClick = {
      "url" : document.URL,
      "href" : href,
      "hrefTop": rect.top,
      "hrefRight": rect.right,
      "hrefBottom": rect.bottom,
      "hrefLeft": rect.left
    }

    jQuery.ajax({
      type: "POST",
      url: "http://192.168.122.150:8080/href-click/",
      data: JSON.stringify(hrefClick),
      dataType: "json",
      contentType: "application/json; charset=utf-8",
      complete: function () {
        window.location.href = href;
      }
    });
  });
});
