/////
// Timeline Data taken from the following match:
// http://matchhistory.na.leagueoflegends.com/en/#match-details/NA1/1653950767
/////

var coord = document.getElementById("coord")
var cords = [JSON.parse(coord.value)],
    // Domain for the current Summoner's Rift on the in-game mini-map
    domain = {
            min: {x: -120, y: -120},
            max: {x: 14870, y: 14980}
    },
    width = 512,
    height = 512,
    bg = "/static/map.png",
    xScale, yScale, svg;

color = d3.scale.linear()
    .domain([0, 3])
    .range(["white", "steelblue"])
    .interpolate(d3.interpolateLab);

xScale = d3.scale.linear()
  .domain([domain.min.x, domain.max.x])
  .range([0, width]);

yScale = d3.scale.linear()
  .domain([domain.min.y, domain.max.y])
  .range([height, 0]);

svg = d3.select("#map").append("svg:svg")
    .attr("width", width)
    .attr("height", height);

svg.append('image')
    .attr('xlink:href', bg)
    .attr('x', '0')
    .attr('y', '0')
    .attr('width', width)
    .attr('height', height);

svg.append('svg:g').selectAll("circle")
    .data(cords)
    .enter().append("svg:circle")
        .attr('cx', function(d) { return xScale(d[0]) })
        .attr('cy', function(d) { return yScale(d[1]) })
        .attr('r', 5)
        .style("fill", "red")
        .attr('class', 'avg');


txt = d3.select("#out").append("p")
    .attr('x', '550')
    .attr('y', '550')
    .text("x,y = " + coord.value);
