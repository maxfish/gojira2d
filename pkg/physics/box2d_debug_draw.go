package physics

import (
	"github.com/ByteArena/box2d"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/maxfish/gojira2d/pkg/graphics"
	"github.com/maxfish/gojira2d/pkg/utils"
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

const numSegmentsPerCircle = 12

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
		circlePoints, _ := utils.CircleToPolygon(mgl64.Vec2{circle.M_p.X, circle.M_p.Y}, circle.M_radius, numSegmentsPerCircle, 0)
		var vertices []mgl64.Vec2
		for i := 0; i < len(circlePoints); i++ {
			vertices = append(vertices, mgl64.Vec2{circlePoints[i].X(), circlePoints[i].Y()})
		}
		c := graphics.NewPolylinePrimitive(mgl64.Vec3{body.GetPosition().X * d.PTM, body.GetPosition().Y * d.PTM, 0}, vertices, true)
		c.SetScale(mgl64.Vec2{d.PTM, d.PTM})
		fixture.SetUserData(c)
	case box2d.B2Shape_Type.E_polygon:
		b2Shape := fixture.GetShape().(*box2d.B2PolygonShape)
		numVertices := b2Shape.M_count
		var vertices []mgl64.Vec2
		for i := 0; i < numVertices; i++ {
			vertices = append(vertices, mgl64.Vec2{b2Shape.M_vertices[i].X, b2Shape.M_vertices[i].Y})
		}
		c := graphics.NewPolylinePrimitive(mgl64.Vec3{body.GetPosition().X * d.PTM, body.GetPosition().Y * d.PTM, 0}, vertices, true)
		c.SetScale(mgl64.Vec2{d.PTM, d.PTM})
		fixture.SetUserData(c)
	case box2d.B2Shape_Type.E_chain:
		b2Shape := fixture.GetShape().(*box2d.B2ChainShape)
		numVertices := b2Shape.M_count
		var vertices []mgl64.Vec2
		for i := 0; i < numVertices; i++ {
			vertices = append(vertices, mgl64.Vec2{b2Shape.M_vertices[i].X, b2Shape.M_vertices[i].Y})
		}
		c := graphics.NewPolylinePrimitive(mgl64.Vec3{body.GetPosition().X * d.PTM, body.GetPosition().Y * d.PTM, 0}, vertices, false)
		c.SetScale(mgl64.Vec2{d.PTM, d.PTM})
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
	case box2d.B2Shape_Type.E_circle, box2d.B2Shape_Type.E_polygon:
		c := fixture.GetUserData().(*graphics.Primitive2D)
		c.SetPosition(mgl64.Vec3{body.GetPosition().X * d.PTM, body.GetPosition().Y * d.PTM, 0})
		c.SetAngle(body.GetAngle())
		c.SetColor(color)
	case box2d.B2Shape_Type.E_chain:
		c := fixture.GetUserData().(*graphics.Primitive2D)
		c.SetPosition(mgl64.Vec3{body.GetPosition().X * d.PTM, body.GetPosition().Y * d.PTM, 0})
		c.SetAngle(body.GetAngle())
		c.SetColor(color)
		// TODO: Chain is not complete yet
	default:
		// There are no other shapes supported
	}
}
