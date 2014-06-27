var autoCompleter = new Bloodhound({
  datumTokenizer: Bloodhound.tokenizers.obj.whitespace('value'),
  queryTokenizer: Bloodhound.tokenizers.whitespace,
  remote: {
    url: '/movies/%QUERY',
    filter: function(list) {
      return $.map(list, function(title) { return { title: title }; });
    }
  }
});

var movie1 = {title: "", tags: [], unique_tags: []};
var movie2 = {title: "", tags: [], unique_tags: []};
var movies = [movie1, movie2];

$(function() {
  autoCompleter.initialize();

  selector = 'input.typeahead';
  bindAutocomplete(selector);
  bindSelectEvent(selector);

  var movieA = getUrlVar("a");
  var movieB = getUrlVar("b");
  if (movieA != "" && movieB != "") {
    $('#movie1 .typeahead').typeahead('val', movieA);
    $('#movie2 .typeahead').typeahead('val', movieB);
    movie1.title = movieA;
    movie2.title = movieB;
    createTags(movie1);
    createTags(movie2);
  }
});

function bindAutocomplete(selector) {
  $(selector).typeahead(null, {
    displayKey: 'title',
    source: autoCompleter.ttAdapter()
  })
}

function bindSelectEvent(selector) {
  $(selector).on("typeahead:selected", function(e, suggestion) {
    var order = e.target.dataset["order"]; // 1 or 2
    movies[order - 1].title = suggestion.title;
    createTags(movies[order - 1]);
  })
}

function createTags(movie) {
  $.get("/tags/" + movie.title, function(tags) {
    movie.tags = tags;
    if (areAllTagsReceived()) {
      populateVenn();
      setUrl();
    }
  })
}

function areAllTagsReceived() {
  for (var j = 0; j < movies.length; j++) {
    if (movies[j].tags.length == 0) {
      return false;
    }
  }
  return true;
}

function populateVenn() {
  commonTags = _.intersection(movie1.tags, movie2.tags);
  movie1.uniqueTags = _.difference(movie1.tags, commonTags);
  movie2.uniqueTags = _.difference(movie2.tags, commonTags);

  $(".common-tags").html(tagsToHtml(commonTags));
  $(".movie1-tags").html(tagsToHtml(movie1.uniqueTags));
  $(".movie2-tags").html(tagsToHtml(movie2.uniqueTags));
  $(".venn").show();
}

function setUrl() {
  var params = "?a=" + encodeURIComponent(movie1.title) + "&b=" + encodeURIComponent(movie2.title);
  history.pushState(null, null, window.location.protocol + "//" + window.location.host + params);
}

function tagsToHtml(tags) {
  var html = "";
  for (var i = 0; i < tags.length; i++) {
    html += '<span><a target="_blank" href=http://www.imdb.com/keyword/' + tags[i] + '>' + tags[i] + "</a></span>";
  }
  return html;
}

function getUrlVar(key){
  var result = new RegExp(key + "=([^&]*)", "i").exec(window.location.search);
  return result && unescape(result[1]) || "";
}
