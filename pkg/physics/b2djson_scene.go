package physics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ByteArena/box2d"
)

type B2DJsonScene struct {
	World       *box2d.B2World
	loadedData  B2DJsonWorld
	indexToBody map[int]*box2d.B2Body
	bodyToName  map[*box2d.B2Body]string

	// Engine parameters
	PositionIterations int
	VelocityIterations int
	StepsPerSecond     float64
}

func NewB2DJsonSceneFromFile(fileName string) *B2DJsonScene {
	scene := &B2DJsonScene{}
	scene.indexToBody = make(map[int]*box2d.B2Body)
	scene.bodyToName = make(map[*box2d.B2Body]string)

	// Open the json file
	jsonFile, err := os.Open(fileName)
	defer jsonFile.Close()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Successfully Opened json")
	}

	// Read the data into the scene structure
	byteData, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteData, &scene.loadedData)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Data read successfully")
	}

	// Engine parameters
	scene.PositionIterations = scene.loadedData.PositionIterations
	scene.VelocityIterations = scene.loadedData.VelocityIterations
	scene.StepsPerSecond = scene.loadedData.StepsPerSecond

	// Build the Box2D objects
	scene.World = scene.buildWorld()
	scene.loadWorld()

	return scene
}

func (s *B2DJsonScene) buildWorld() *box2d.B2World {
	w := s.loadedData
	b2World := box2d.MakeB2World(box2d.MakeB2Vec2(w.Gravity.X, w.Gravity.Y))
	b2World.SetAllowSleeping(w.AllowSleep)
	b2World.SetAutoClearForces(w.AutoClearForces)
	b2World.M_warmStarting = w.WarmStarting
	b2World.M_continuousPhysics = w.ContinuousPhysics
	b2World.M_subStepping = w.SubStepping

	return &b2World
}

func (s *B2DJsonScene) loadWorld() {
	w := s.loadedData

	for i := 0; i < len(w.Body); i++ {
		bodyData := w.Body[i]
		body := s.buildBody(&bodyData)
		s.indexToBody[i] = body
	}

	//	//need two passes for joints because gear joints reference other joints
	//	i = 0;
	//Json::Value jointValue = worldValue["joint"][i++];
	//	while ( !jointValue.isNull() ) {
	//		if ( jointValue["type"].asString() != "gear" ) {
	//			b2Joint* joint = j2b2Joint(world, jointValue);
	//			readCustomPropertiesFromJson(joint, jointValue);
	//			m_joints.push_back(joint);
	//		}
	//		jointValue = worldValue["joint"][i++];
	//	}
	//	i = 0;
	//	jointValue = worldValue["joint"][i++];
	//	while ( !jointValue.isNull() ) {
	//		if ( jointValue["type"].asString() == "gear" ) {
	//			b2Joint* joint = j2b2Joint(world, jointValue);
	//			readCustomPropertiesFromJson(joint, jointValue);
	//			m_joints.push_back(joint);
	//		}
	//		jointValue = worldValue["joint"][i++];
	//	}
	//
	//	i = 0;
	//Json::Value imageValue = worldValue["image"][i++];
	//	while ( !imageValue.isNull() ) {
	//		b2dJsonImage* img = j2b2dJsonImage(imageValue);
	//		readCustomPropertiesFromJson(img, imageValue);
	//		m_images.push_back(img);
	//		addImage(img);
	//
	//		imageValue = worldValue["image"][i++];
	//	}
	//
	//	return world;
}

func (s *B2DJsonScene) buildBody(data *B2DBodyData) *box2d.B2Body {
	var b2BodyDef = box2d.MakeB2BodyDef()
	b2BodyDef.Type = data.Type
	b2BodyDef.Position = box2d.B2Vec2{X: data.Position.X, Y: data.Position.Y}
	b2BodyDef.Angle = data.Angle
	b2BodyDef.LinearVelocity = box2d.B2Vec2{X: data.LinearVelocity.X, Y: data.LinearVelocity.Y}
	b2BodyDef.AngularVelocity = data.AngularVelocity
	b2BodyDef.LinearDamping = data.LinearDamping
	b2BodyDef.AngularDamping = data.AngularDamping
	b2BodyDef.Awake = data.Awake
	b2BodyDef.FixedRotation = data.FixedRotation
	b2BodyDef.Bullet = data.Bullet
	//b2BodyDef.GravityScale = 1 // Value not loaded from file
	//b2BodyDef.AllowSleep = true // Value not loaded from file
	//b2BodyDef.Active = true // Value not loaded from file

	b2Body := s.World.CreateBody(&b2BodyDef)
	if data.Name != "" {
		s.bodyToName[b2Body] = data.Name
	}

	for i := 0; i < len(data.Fixture); i++ {
		_ = s.buildFixture(b2Body, &data.Fixture[i])
		////readCustomPropertiesFromJson(fixture, fixtureValue);
	}

	b2MassData := box2d.MakeMassData()
	b2MassData.Mass = data.MassDataMass
	b2MassData.Center = box2d.B2Vec2{X: data.MassDataCenter.X, Y: data.MassDataCenter.Y}
	b2MassData.I = data.MassDataI
	b2Body.SetMassData(&b2MassData)

	return b2Body
}

func (s *B2DJsonScene) buildFixture(b2Body *box2d.B2Body, data *B2DFixtureData) *box2d.B2Fixture {
	var b2Fixture *box2d.B2Fixture
	b2FixtureDef := box2d.MakeB2FixtureDef()
	b2FixtureDef.Restitution = data.Restitution
	b2FixtureDef.Friction = data.Friction
	b2FixtureDef.Density = data.Density
	b2FixtureDef.IsSensor = data.Sensor

	filter := box2d.MakeB2Filter()
	if data.FilterCategoryBits != nil {
		filter.CategoryBits = *data.FilterCategoryBits
	}
	if data.FilterMaskBits != nil {
		filter.MaskBits = *data.FilterMaskBits
	}
	filter.GroupIndex = data.FilterGroupIndex
	b2FixtureDef.Filter = filter

	if data.Circle != nil {
		b2CircleShape := box2d.MakeB2CircleShape()
		b2CircleShape.M_radius = data.Circle.Radius
		b2CircleShape.M_p = box2d.B2Vec2{X: data.Circle.Center.X, Y: data.Circle.Center.Y}
		b2FixtureDef.Shape = b2CircleShape
		b2Fixture = b2Body.CreateFixtureFromDef(&b2FixtureDef)
	} else if data.Polygon != nil {
		var vertices []box2d.B2Vec2
		numVertices := len(data.Polygon.Vertices.X)
		if numVertices > box2d.B2_maxPolygonVertices {
			fmt.Println("Warning: ignoring fixture with too many vertices")
		} else if numVertices < 2 {
			fmt.Println("Warning: ignoring fixture with less than two vertices")
		} else if numVertices == 2 {
			fmt.Println("Warning: creating edge shape instead of polygon with two vertices")
			b2EdgeShape := box2d.MakeB2EdgeShape()
			b2EdgeShape.M_vertex1 = box2d.B2Vec2{X: data.Polygon.Vertices.X[0], Y: data.Polygon.Vertices.Y[0]}
			b2EdgeShape.M_vertex2 = box2d.B2Vec2{X: data.Polygon.Vertices.X[1], Y: data.Polygon.Vertices.Y[1]}
			b2FixtureDef.Shape = &b2EdgeShape
			b2Fixture = b2Body.CreateFixtureFromDef(&b2FixtureDef)
		} else {
			b2PolygonShape := box2d.MakeB2PolygonShape()
			for i := 0; i < numVertices; i++ {
				vertices = append(vertices, box2d.B2Vec2{X: data.Polygon.Vertices.X[i], Y: data.Polygon.Vertices.Y[i]})
			}
			b2PolygonShape.Set(vertices, numVertices)
			b2FixtureDef.Shape = &b2PolygonShape
			b2Fixture = b2Body.CreateFixtureFromDef(&b2FixtureDef)
		}
	} else if data.Chain != nil {
		var vertices []box2d.B2Vec2
		b2ChainShape := box2d.MakeB2ChainShape()
		numVertices := len(data.Chain.Vertices.X)
		for i := 0; i < numVertices; i++ {
			vertices = append(vertices, box2d.B2Vec2{X: data.Chain.Vertices.X[i], Y: data.Chain.Vertices.Y[i]})
		}
		b2ChainShape.CreateChain(vertices, numVertices)
		b2ChainShape.M_hasPrevVertex = data.Chain.HasPrevVertex
		if b2ChainShape.M_hasPrevVertex {
			b2ChainShape.M_prevVertex = box2d.B2Vec2{X: data.Chain.PrevVertex.X, Y: data.Chain.PrevVertex.Y}
		}
		b2ChainShape.M_hasNextVertex = data.Chain.HasNextVertex
		if b2ChainShape.M_hasNextVertex {
			b2ChainShape.M_nextVertex = box2d.B2Vec2{X: data.Chain.NextVertex.X, Y: data.Chain.NextVertex.Y}
		}
		b2FixtureDef.Shape = &b2ChainShape
		b2Fixture = b2Body.CreateFixtureFromDef(&b2FixtureDef)
	}

	//else if ( !fixtureValue["edge"].isNull() ) {
	//	b2EdgeShape edgeShape;
	//	edgeShape.m_vertex1 = jsonToVec("vertex1", fixtureValue["edge"]);
	//	edgeShape.m_vertex2 = jsonToVec("vertex2", fixtureValue["edge"]);
	//	edgeShape.m_hasVertex0 = fixtureValue["edge"].get("hasVertex0",false).asBool();
	//	edgeShape.m_hasVertex3 = fixtureValue["edge"].get("hasVertex3",false).asBool();
	//	if ( edgeShape.m_hasVertex0 )
	//	edgeShape.m_vertex0 = jsonToVec("vertex0", fixtureValue["edge"]);
	//	if ( edgeShape.m_hasVertex3 )
	//	edgeShape.m_vertex3 = jsonToVec("vertex3", fixtureValue["edge"]);
	//	fixtureDef.shape = &edgeShape;
	//	fixture = b2Body->CreateFixture(&fixtureDef);
	//}
	//else if ( !fixtureValue["loop"].isNull() ) { //support old format (r197)
	//	b2ChainShape chainShape;
	//	int numVertices = fixtureValue["loop"]["vertices"]["x"].size();
	//	b2Vec2* vertices = new b2Vec2[numVertices];
	//	for (int i = 0; i < numVertices; i++)
	//	vertices[i] = jsonToVec("vertices", fixtureValue["loop"], i);
	//	chainShape.CreateLoop(vertices, numVertices);
	//	fixtureDef.shape = &chainShape;
	//	fixture = b2Body->CreateFixture(&fixtureDef);
	//	delete[] vertices;
	//}
	//
	//string fixtureName = fixtureValue.get("name","").asString();
	//if ( fixtureName != "" ) {
	//	setFixtureName(fixture, fixtureName.c_str());
	//}
	//
	//string fixturePath = fixtureValue.get("path","").asString();
	//if ( fixturePath != "" ) {
	//	setFixturePath(fixture, fixturePath.c_str());
	//}
	//
	//return fixture;

	return b2Fixture
}
