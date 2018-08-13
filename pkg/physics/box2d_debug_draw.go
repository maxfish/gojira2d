package physics

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/maxfish/box2d"
	"github.com/maxfish/gojira2d/pkg/graphics"
)

type Box2DDebugDraw struct {
	colorNormal    graphics.Color
	colorInactive  graphics.Color
	colorStatic    graphics.Color
	colorKinematic graphics.Color
	colorAsleep    graphics.Color
	colorSensor    graphics.Color

	b2World *box2d.B2World
	PTM     float64
}

func NewBox2DDebugDraw(w *box2d.B2World, PTM float64) *Box2DDebugDraw {
	d := &Box2DDebugDraw{}
	d.b2World = w
	d.PTM = PTM

	d.colorNormal = graphics.Color{0.9, 0.7, 0.7, 1}
	d.colorInactive = graphics.Color{0.5, 0.5, 0.3, 1}
	d.colorStatic = graphics.Color{0.5, 0.9, 0.5, 1}
	d.colorKinematic = graphics.Color{0.5, 0.5, 0.9, 1}
	d.colorAsleep = graphics.Color{0.6, 0.6, 0.6, 1}
	d.colorSensor = graphics.Color{0.6, 0.3, 0.6, 1}

	// TODO: Debug flags??
	drawShapes := true

	if drawShapes {
		body := w.GetBodyList()
		for body != nil {
			fixture := body.GetFixtureList()
			for fixture != nil {
				d.buildShape(body, fixture)
				fixture = fixture.GetNext()
			}
			body = body.GetNext()
		}
	}

	// TODO Draw joints

	return d
}

func (d *Box2DDebugDraw) Update() {
	color := graphics.Color{}
	body := d.b2World.GetBodyList()
	for body != nil {
		transform := body.GetTransform()
		fixture := body.GetFixtureList()
		for fixture != nil {
			if fixture.IsSensor() {
				color = d.colorSensor
			} else if !body.IsActive() {
				color = d.colorInactive
			} else if body.GetType() == box2d.B2BodyType.B2_staticBody {
				color = d.colorStatic
			} else if body.GetType() == box2d.B2BodyType.B2_kinematicBody {
				color = d.colorKinematic
			} else if !body.IsAwake() {
				color = d.colorAsleep
			} else {
				color = d.colorNormal
			}
			d.updateShape(body, fixture, transform, color)
			fixture = fixture.GetNext()
		}
		body = body.GetNext()
	}
}

func (d *Box2DDebugDraw) Draw(context *graphics.Context) {
	body := d.b2World.GetBodyList()
	for body != nil {
		fixture := body.GetFixtureList()
		for fixture != nil {
			primitive := fixture.GetUserData().(*graphics.Primitive2D)
			primitive.Draw(context)
			fixture = fixture.GetNext()
		}
		body = body.GetNext()
	}
}

func (d *Box2DDebugDraw) buildShape(body *box2d.B2Body, fixture *box2d.B2Fixture) {
	switch fixture.GetType() {
	case box2d.B2Shape_Type.E_circle:
		circle := fixture.GetShape().(*box2d.B2CircleShape)
		c := graphics.NewRegularPolygonPrimitive(mgl32.Vec3{float32(body.GetPosition().X * d.PTM), float32(body.GetPosition().Y * d.PTM), 0}, float32(circle.M_radius*d.PTM), 10, false)
		c.SetAnchorToCenter()
		fixture.SetUserData(c)
	case box2d.B2Shape_Type.E_polygon:
		b2Shape := fixture.GetShape().(*box2d.B2PolygonShape)
		numVertices := b2Shape.M_count
		var vertices []mgl32.Vec2
		for i := 0; i < numVertices; i++ {
			vertices = append(vertices, mgl32.Vec2{float32(b2Shape.M_vertices[i].X), float32(b2Shape.M_vertices[i].Y)})
		}
		c := graphics.NewPolylinePrimitiveRaw(mgl32.Vec3{float32(body.GetPosition().X * d.PTM), float32(body.GetPosition().Y * d.PTM), 0}, vertices, true)
		c.SetScale(mgl32.Vec2{float32(d.PTM), float32(d.PTM)})
		fixture.SetUserData(c)
	case box2d.B2Shape_Type.E_chain:
		b2Shape := fixture.GetShape().(*box2d.B2ChainShape)
		numVertices := b2Shape.M_count
		var vertices []mgl32.Vec2
		for i := 0; i < numVertices; i++ {
			vertices = append(vertices, mgl32.Vec2{float32(b2Shape.M_vertices[i].X), float32(b2Shape.M_vertices[i].Y)})
		}
		c := graphics.NewPolylinePrimitiveRaw(mgl32.Vec3{float32(body.GetPosition().X * d.PTM), float32(body.GetPosition().Y * d.PTM), 0}, vertices, false)
		c.SetScale(mgl32.Vec2{float32(d.PTM), float32(d.PTM)})
		fixture.SetUserData(c)

		// 			if (chain.m_hasPrevVertex)
		// 			{
		// 				b2Vec2 vp = b2Mul(xf, chain.m_prevVertex);
		// 				g_debugDraw.DrawSegment(vp, v1, ghostColor);
		// 				g_debugDraw.DrawCircle(vp, 0.1f, ghostColor);
		// 			}

		// 			for (int i = 1; i < count; ++i)
		// 			{
		// 				b2Vec2 v2 = b2Mul(xf, vertices[i]);
		// 				g_debugDraw.DrawSegment(v1, v2, color);
		// 				g_debugDraw.DrawPoint(v2, 4.0, color);
		// 				v1 = v2;
		// 			}

		// 			if (chain.m_hasNextVertex)
		// 			{
		// 				b2Vec2 vn = b2Mul(xf, chain.m_nextVertex);
		// 				g_debugDraw.DrawSegment(v1, vn, ghostColor);
		// 				g_debugDraw.DrawCircle(vn, 0.1f, ghostColor);
		// 			}
	default:
		// There are no other shapes supported
	}
}

func (d *Box2DDebugDraw) updateShape(body *box2d.B2Body, fixture *box2d.B2Fixture, transform box2d.B2Transform, color graphics.Color) {
	switch fixture.GetType() {
	case box2d.B2Shape_Type.E_circle:
		c := fixture.GetUserData().(*graphics.Primitive2D)
		c.SetPosition(mgl32.Vec3{float32(body.GetPosition().X * d.PTM), float32(body.GetPosition().Y * d.PTM), 0})
		c.SetAngle(float32(body.GetAngle()))
		c.SetColor(color)
	case box2d.B2Shape_Type.E_polygon:
		c := fixture.GetUserData().(*graphics.Primitive2D)
		c.SetPosition(mgl32.Vec3{float32(body.GetPosition().X * d.PTM), float32(body.GetPosition().Y * d.PTM), 0})
		c.SetAngle(float32(body.GetAngle()))
		c.SetColor(color)
	case box2d.B2Shape_Type.E_chain:
		c := fixture.GetUserData().(*graphics.Primitive2D)
		c.SetPosition(mgl32.Vec3{float32(body.GetPosition().X * d.PTM), float32(body.GetPosition().Y * d.PTM), 0})
		c.SetAngle(float32(body.GetAngle()))
		c.SetColor(color)
		// TODO: Chain is not complete yet
	default:
		// There are no other shapes supported
	}
}
