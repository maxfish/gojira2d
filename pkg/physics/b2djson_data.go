package physics

import "encoding/json"

// Description of a scene as produced by R.U.B.E editor (https://www.iforce2d.net/rube/json-structure)

type B2DJsonWorld struct {
	Gravity            B2DVector2D               `json:"gravity"`
	AllowSleep         bool                      `json:"allowSleep"`
	AutoClearForces    bool                      `json:"autoClearForces"`
	PositionIterations int                       `json:"positionIterations"`
	VelocityIterations int                       `json:"velocityIterations"`
	StepsPerSecond     float64                   `json:"stepsPerSecond"`
	SubStepping        bool                      `json:"subStepping"`
	WarmStarting       bool                      `json:"warmStarting"`
	ContinuousPhysics  bool                      `json:"continuousPhysics"`
	Collisionbitplanes B2DCollisionBitplanesData `json:"collisionbitplanes"`

	Body  []B2DBodyData  `json:"body"`
	Image []B2DImageData `json:"image"`
	Joint []B2DJointData `json:"joint"`
}

type B2DBodyData struct {
	Name            string      `json:"name"`
	Type            uint8       `json:"type"`  // 0 = static, 1 = kinematic, 2 = dynamic
	Angle           float64     `json:"angle"` // radians
	AngularDamping  float64     `json:"angularDamping"`
	AngularVelocity float64     `json:"angularVelocity"` // radians per second
	Awake           bool        `json:"awake"`
	Bullet          bool        `json:"bullet"`
	FixedRotation   bool        `json:"fixedRotation"`
	LinearVelocity  B2DVector2D `json:"linearVelocity"`
	LinearDamping   float64     `json:"linearDamping"`
	MassDataI       float64     `json:"massData-I"`
	MassDataMass    float64     `json:"massData-mass"`
	MassDataCenter  B2DVector2D `json:"massData-center"`
	Position        B2DVector2D `json:"position"`

	Fixture          []B2DFixtureData        `json:"fixture"`
	CustomProperties []B2DCustomPropertyData `json:"customProperties"`
}

type B2DFixtureData struct {
	Name               string  `json:"name"`
	Density            float64 `json:"density"`
	FilterCategoryBits *uint16 `json:"filter-categoryBits"` // if not present, interpret as 1
	FilterMaskBits     *uint16 `json:"filter-maskBits"`     // if not present, interpret as 65535
	FilterGroupIndex   int16   `json:"filter-groupIndex"`
	Friction           float64 `json:"friction"`
	Restitution        float64 `json:"restitution"`
	Sensor             bool    `json:"sensor"`
	// A fixture object will have only one of the following shape objects
	Circle  *B2DCircleFixtureData  `json:"circle"`
	Polygon *B2DPolygonFixtureData `json:"polygon"`
	Chain   *B2DChainFixtureData   `json:"chain"`

	CustomProperties []B2DCustomPropertyData `json:"customProperties"`
}

type B2DJointData struct {
	Type             string      `json:"type"`
	Name             string      `json:"name"`
	AnchorA          B2DVector2D `json:"anchorA"`
	AnchorB          B2DVector2D `json:"anchorB"`
	BodyA            int         `json:"bodyA"` // zero-based index of body in bodies array
	BodyB            int         `json:"bodyB"` // zero-based index of body in bodies array
	CollideConnected bool        `json:"collideConnected"`

	DampingRatio       float64     `json:"dampingRatio"`       // Distance, Weld
	CorrectionFactor   float64     `json:"correctionFactor"`   // Motor
	EnableLimit        bool        `json:"enableLimit"`        // Revolute, Prismatic
	EnableMotor        bool        `json:"enableMotor"`        // Revolute, Prismatic, Wheel
	Frequency          float64     `json:"frequency"`          // Distance, Weld
	JointSpeed         float64     `json:"jointSpeed"`         // Revolute
	Length             float64     `json:"length"`             // Distance
	LocalAxisA         B2DVector2D `json:"localAxisA"`         // Prismatic, Wheel
	LowerLimit         float64     `json:"lowerLimit"`         // Revolute, Prismatic
	MaxForce           float64     `json:"maxForce"`           // Motor, Friction
	MaxLength          float64     `json:"maxLength"`          // Rope
	MaxMotorForce      float64     `json:"maxMotorForce"`      // Prismatic
	MaxMotorTorque     float64     `json:"maxMotorTorque"`     // Revolute, Wheel
	MaxTorque          float64     `json:"maxTorque"`          // Motor, Friction
	MotorSpeed         float64     `json:"motorSpeed"`         // Revolute, Prismatic, Wheel
	RefAngle           float64     `json:"refAngle"`           // Revolute, Prismatic, Weld
	SpringDampingRatio float64     `json:"springDampingRatio"` // Wheel
	SpringFrequency    float64     `json:"springFrequency"`    // Wheel
	UpperLimit         float64     `json:"upperLimit"`         // Revolute, Prismatic

	CustomProperties []B2DCustomPropertyData `json:"customProperties"`
}

type B2DImageData struct {
	Name        string          `json:"name"`
	Opacity     float64         `json:"opacity"`
	RenderOrder float64         `json:"renderOrder"`
	Scale       float64         `json:"scale"`       // the length of the vertical side of the image in physics units
	AspectScale float64         `json:"aspectScale"` // the ratio of width to height, relative to the natural dimensions
	Angle       float64         `json:"angle"`       // radians
	Body        int             `json:"body"`        // zero-based index of body in bodies array
	Center      B2DVector2D     `json:"center"`      // center position in body local coordinates
	Corners     B2DVerticesData `json:"corners"`     // corner positions in body local coordinates
	File        string          `json:"file"`        // if relative, from the location of the exported file
	Filter      int             `json:"filter"`      // texture magnification filter, 0 = linear, 1 = nearest
	Flip        bool            `json:"flip"`        // true if the texture should be reversed horizontally

	ColorTint         []int                   `json:"colorTint"`         // RGBA values for color tint, if not 255,255,255,255
	GlDrawElements    []int                   `json:"glDrawElements"`    //Indices for drawing GL_TRIANGLES with the glDrawElements function and the other glXXX properties below
	GlTexCoordPointer []float64               `json:"glTexCoordPointer"` //Texture coordinates for use with glTexCoordPointer (the 'flip' property has already been taken into account)
	GlVertexPointer   []float64               `json:"glVertexPointer"`   // Vertex positions for use with glVertexPointer
	CustomProperties  []B2DCustomPropertyData `json:"customProperties"`
}

type B2DVerticesData struct {
	X []float64 `json:"x"`
	Y []float64 `json:"y"`
}

type B2DCollisionBitplanesData struct {
	Names []string `json:"names"`
}

type B2DPolygonFixtureData struct {
	Vertices B2DVerticesData `json:"vertices"`
}

type B2DChainFixtureData struct {
	Vertices B2DVerticesData `json:"vertices"`
	//If the following properties are not present, the shape is an open-ended
	//chain shape. If they are present, the shape is a closed loop shape.
	HasNextVertex bool        `json:"hasNextVertex"`
	HasPrevVertex bool        `json:"hasPrevVertex"`
	NextVertex    B2DVector2D `json:"nextVertex"`
	PrevVertex    B2DVector2D `json:"prevVertex"`
}

type B2DCircleFixtureData struct {
	Center B2DVector2D `json:"center"`
	Radius float64     `json:"radius"`
}

type B2DVector2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Unmarshaler for the Vector2D type.
type _vector2D B2DVector2D

func (a *B2DVector2D) UnmarshalJSON(b []byte) (err error) {
	j, n := _vector2D{}, float64(0)
	if err = json.Unmarshal(b, &j); err == nil {
		*a = B2DVector2D(j)
		return
	}
	if err = json.Unmarshal(b, &n); err == nil {
		a.X = n
		a.Y = n
	}
	return
}

type B2DCustomPropertyData struct {
	// TODO Implement
}
