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

var movie1 = {tags: [], unique_tags: []};
var movie2 = {tags: [], unique_tags: []};
var movies = [movie1, movie2];

$(function() {
  autoCompleter.initialize();

  selector = 'input.typeahead';
  bindAutocomplete(selector);
  bindSelectEvent(selector);
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
    $.get("/tags/" + suggestion.title, function(tags) {
      movies[order - 1].tags = tags;

      if (areAllTagsReceived()) {
        populateVenn();
      }
    })
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
}

function tagsToHtml(tags) {
  var html = "";
  for (var i = 0; i < tags.length; i++) {
    html += "<span>" + tags[i] + "</span>";
  }
  return html;
}
