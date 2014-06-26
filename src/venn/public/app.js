$(function() {
  var movie = new Bloodhound({
    datumTokenizer: Bloodhound.tokenizers.obj.whitespace('value'),
    queryTokenizer: Bloodhound.tokenizers.whitespace,
    remote: {
      url: '/movies/%QUERY',
      filter: function(list) {
        return $.map(list, function(title) { return { title: title }; });
      }
    }
  });

  movie.initialize();

  $('#movie1 .typeahead').typeahead(null, {
    displayKey: 'title',
    source: movie.ttAdapter()
  });
});
