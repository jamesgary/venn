var chosenKeywords = [];

$(function() {
  populateNewGame()
});

function populateNewGame() {
  $.get("/new_game", function(data) {
    var keyword = data["keyword"];
    var movieCount = data["movieCount"];
    var possibleKeywords = data["possibleKeywords"];
    chosenKeywords.push(keyword);

    $(".chosen.keywords").html(keywordsToHtml(chosenKeywords));
    $(".choices.keywords").html(keywordsToHtml(possibleKeywords));
    $(".movie-count").html(movieCount);

    bindChoosingKeyword();
  });
}

function bindChoosingKeyword() {
  $(".choices .keyword").click(function(e) {
    chosenKeywords.push(e.target.dataset["keyword"]);
    $(".chosen.keywords").html(keywordsToHtml(chosenKeywords));
    $.get("/guess/" + chosenKeywords.join(" "), function(data) {
      if (data["movie"]) {
        win(data["movie"]);
      } else {
        var movieCount = data["movieCount"];
        var possibleKeywords = data["possibleKeywords"];
        $(".choices.keywords").html(keywordsToHtml(possibleKeywords));
        $(".movie-count").html(movieCount);
        bindChoosingKeyword();
      }
    })
  })
}

function win(movie) {
  $(".status").html('You win! The only movie in the universe to match these keywords is: <a target="_blank" href="http://www.imdb.com/keyword/' + chosenKeywords.join('/') + '">' + movie + '</a>');
  $(".in-game").hide();
  $(".restart").show();
}

function keywordsToHtml(keywords) {
  var html = "";
  for (var i = 0; i < keywords.length; i++) {
    html += '<span class="keyword" data-keyword="' + keywords[i] + '">' + keywords[i] + "</a></span>";
  }
  return html;
}

function getUrlVar(key){
  var result = new RegExp(key + "=([^&]*)", "i").exec(window.location.search);
  return result && unescape(result[1]) || "";
}
