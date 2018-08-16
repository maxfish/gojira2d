package physics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/maxfish/box2d"
)

// B2DJsonScene A scene produced by R.U.B.E
type B2DJsonScene struct {
	World       *box2d.B2World
	loadedData  B2DJsonWorld
	indexToBody map[int]*box2d.B2Body
	bodyToName  map[*box2d.B2Body]string
	nameToBody  map[string]*box2d.B2Body
	jointToName map[*box2d.B2Joint]string
	bodies      []*box2d.B2Body
	joints      []*box2d.B2Joint

	// Engine parameters
	PositionIterations int
	VelocityIterations int
	StepsPerSecond     float64
}

// NewB2DJsonSceneFromFile Loads the scene from an exported JSON file
func NewB2DJsonSceneFromFile(fileName string) *B2DJsonScene {
	scene := &B2DJsonScene{}
	scene.indexToBody = make(map[int]*box2d.B2Body)
	scene.bodyToName = make(map[*box2d.B2Body]string)
	scene.nameToBody = make(map[string]*box2d.B2Body)
	scene.jointToName = make(map[*box2d.B2Joint]string)

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
		fmt.Println(err)
	} else {
		fmt.Println("Data read successfully")
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

func (s *B2DJsonScene) BodyForName(name string) *box2d.B2Body {
	return s.nameToBody[name]
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
		s.nameToBody[bodyData.Name] = body
		s.bodies = append(s.bodies, body)
	}

	// NOTE: R.U.B.E doesn't support Gear joints. To support them in this loader
	// two loops are needed. The first one should parse all the non-gear joints and the second
	// only the gears. Gear joints reference other joins.
	for i := 0; i < len(w.Joint); i++ {
		jointData := w.Joint[i]
		_ = s.buildJoint(&jointData)
		//s.joints = append(s.joints, joint)
	}
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

	b2Body := s.World.CreateBody(&b2BodyDef)
	if data.Name != "" {
		s.bodyToName[b2Body] = data.Name
	}

	for i := 0; i < len(data.Fixture); i++ {
		_ = s.buildFixture(b2Body, &data.Fixture[i])
		// TODO: Read the custom properties
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

	// TODO: Read the name of the fixture and store it

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

	// TODO: Read the custom properties
	// NOTE: these shapes are not exported by R.U.B.E -> edge, loop

	return b2Fixture
}

func (s *B2DJsonScene) buildJoint(data *B2DJointData) box2d.B2JointInterface {
	var jointInterface box2d.B2JointInterface
	bodyIndexA := data.BodyA
	bodyIndexB := data.BodyB

	if bodyIndexA >= len(s.bodies) || bodyIndexB >= len(s.bodies) {
		fmt.Println("Error: couldn't create the joint. Bodies indices wrong")
		return nil
	}

	// TODO: Read the name of the joint and store it

	if data.Type == "revolute" {
		j := box2d.MakeB2RevoluteJointDef()
		j.BodyA = s.bodies[bodyIndexA]
		j.BodyB = s.bodies[bodyIndexB]
		j.CollideConnected = data.CollideConnected
		j.LocalAnchorA = box2d.B2Vec2{X: data.AnchorA.X, Y: data.AnchorA.Y}
		j.LocalAnchorB = box2d.B2Vec2{X: data.AnchorB.X, Y: data.AnchorB.Y}
		j.ReferenceAngle = data.RefAngle
		j.EnableLimit = data.EnableLimit
		j.LowerAngle = data.LowerLimit
		j.UpperAngle = data.UpperLimit
		j.EnableMotor = data.EnableMotor
		j.MotorSpeed = data.MotorSpeed
		j.MaxMotorTorque = data.MaxMotorTorque
		jointInterface = s.World.CreateJoint(&j)
	} else if data.Type == "prismatic" {
		// TODO: This works only with fixed rotation set to true!
		s.bodies[bodyIndexA].SetFixedRotation(true)
		s.bodies[bodyIndexB].SetFixedRotation(true)

		j := box2d.MakeB2PrismaticJointDef()
		j.SetBodyA(s.bodies[bodyIndexA])
		j.SetBodyB(s.bodies[bodyIndexB])
		j.SetCollideConnected(data.CollideConnected)
		j.LocalAnchorA = box2d.B2Vec2{X: data.AnchorA.X, Y: data.AnchorA.Y}
		j.LocalAnchorB = box2d.B2Vec2{X: data.AnchorB.X, Y: data.AnchorB.Y}
		j.LocalAxisA = box2d.B2Vec2{X: data.LocalAxisA.X, Y: data.LocalAxisA.Y}
		j.ReferenceAngle = data.RefAngle
		j.EnableLimit = data.EnableLimit
		j.EnableMotor = data.EnableMotor
		j.MotorSpeed = data.MotorSpeed
		j.MaxMotorForce = data.MaxMotorForce
		j.LowerTranslation = data.LowerLimit
		j.UpperTranslation = data.UpperLimit
		jointInterface = s.World.CreateJoint(&j)
	} else if data.Type == "distance" {
		j := box2d.MakeB2DistanceJointDef()
		j.SetBodyA(s.bodies[bodyIndexA])
		j.SetBodyB(s.bodies[bodyIndexB])
		j.SetCollideConnected(data.CollideConnected)
		j.LocalAnchorA = box2d.B2Vec2{X: data.AnchorA.X, Y: data.AnchorA.Y}
		j.LocalAnchorB = box2d.B2Vec2{X: data.AnchorB.X, Y: data.AnchorB.Y}
		j.Length = data.Length
		j.FrequencyHz = data.Frequency
		j.DampingRatio = data.DampingRatio
		jointInterface = s.World.CreateJoint(&j)
	} else if data.Type == "wheel" {
		j := box2d.MakeB2WheelJointDef()
		j.BodyA = s.bodies[bodyIndexA]
		j.BodyB = s.bodies[bodyIndexB]
		j.CollideConnected = data.CollideConnected
		j.LocalAnchorA = box2d.B2Vec2{X: data.AnchorA.X, Y: data.AnchorA.Y}
		j.LocalAnchorB = box2d.B2Vec2{X: data.AnchorB.X, Y: data.AnchorB.Y}
		j.LocalAxisA = box2d.B2Vec2{X: data.LocalAxisA.X, Y: data.LocalAxisA.Y}
		j.EnableMotor = data.EnableMotor
		j.MotorSpeed = data.MotorSpeed
		j.MaxMotorTorque = data.MaxMotorTorque
		j.FrequencyHz = data.SpringFrequency
		j.DampingRatio = data.SpringDampingRatio
		jointInterface = s.World.CreateJoint(&j)
	} else if data.Type == "weld" {
		j := box2d.MakeB2WeldJointDef()
		j.SetBodyA(s.bodies[bodyIndexA])
		j.SetBodyB(s.bodies[bodyIndexB])
		j.SetCollideConnected(data.CollideConnected)
		j.LocalAnchorA = box2d.B2Vec2{X: data.AnchorA.X, Y: data.AnchorA.Y}
		j.LocalAnchorB = box2d.B2Vec2{X: data.AnchorB.X, Y: data.AnchorB.Y}
		j.ReferenceAngle = data.RefAngle
		j.FrequencyHz = data.Frequency
		j.DampingRatio = data.DampingRatio
		jointInterface = s.World.CreateJoint(&j)
	} else if data.Type == "rope" {
		j := box2d.MakeB2RopeJointDef()
		j.SetBodyA(s.bodies[bodyIndexA])
		j.SetBodyB(s.bodies[bodyIndexB])
		j.SetCollideConnected(data.CollideConnected)
		j.LocalAnchorA = box2d.B2Vec2{X: data.AnchorA.X, Y: data.AnchorA.Y}
		j.LocalAnchorB = box2d.B2Vec2{X: data.AnchorB.X, Y: data.AnchorB.Y}
		j.MaxLength = data.MaxLength
		jointInterface = s.World.CreateJoint(&j)
	} else if data.Type == "motor" {
		j := box2d.MakeB2MotorJointDef()
		j.SetBodyA(s.bodies[bodyIndexA])
		j.SetBodyB(s.bodies[bodyIndexB])
		j.SetCollideConnected(data.CollideConnected)
		j.LinearOffset = box2d.B2Vec2{X: data.LinearOffset.X, Y: data.LinearOffset.Y}
		j.AngularOffset = data.RefAngle
		j.MaxForce = data.MaxForce
		j.MaxTorque = data.MaxTorque
		j.CorrectionFactor = data.CorrectionFactor
		jointInterface = s.World.CreateJoint(&j)
	} else if data.Type == "friction" {
		j := box2d.MakeB2FrictionJointDef()
		j.SetBodyA(s.bodies[bodyIndexA])
		j.SetBodyB(s.bodies[bodyIndexB])
		j.SetCollideConnected(data.CollideConnected)
		j.LocalAnchorA = box2d.B2Vec2{X: data.AnchorA.X, Y: data.AnchorA.Y}
		j.LocalAnchorB = box2d.B2Vec2{X: data.AnchorB.X, Y: data.AnchorB.Y}
		j.MaxForce = data.MaxForce
		j.MaxTorque = data.MaxTorque
		jointInterface = s.World.CreateJoint(&j)
	}

	return jointInterface
}
