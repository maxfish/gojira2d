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

	// TODO: Debug flags??
	drawShapes := true

	if drawShapes {
		body := w.GetBodyList()
		for body != nil {
			color := graphics.Color{}
			transform := body.GetTransform()
			fixture := body.GetFixtureList()
			for fixture != nil {
				if !body.IsActive() {
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
				d.buildShape(body, fixture, transform, color)
				fixture = fixture.GetNext()
			}
			body = body.GetNext()
		}
	}

	return d
}

func (d *Box2DDebugDraw) Draw(context *graphics.Context) {
	body := d.b2World.GetBodyList()
	for body != nil {
		if body.GetUserData() != nil {
			primitive := body.GetUserData().(*graphics.Primitive2D)
			primitive.Draw(context)
		}
		body = body.GetNext()
	}
}

func (d *Box2DDebugDraw) Update() {
	body := d.b2World.GetBodyList()
	for body != nil {
		color := graphics.Color{}
		transform := body.GetTransform()
		fixture := body.GetFixtureList()
		for fixture != nil {
			if !body.IsActive() {
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

func (d *Box2DDebugDraw) buildShape(body *box2d.B2Body, fixture *box2d.B2Fixture, transform box2d.B2Transform, color graphics.Color) {
	switch fixture.GetType() {
	case box2d.B2Shape_Type.E_circle:
		circle := fixture.GetShape().(*box2d.B2CircleShape)
		c := graphics.NewRegularPolygonPrimitive(mgl32.Vec3{float32(body.GetPosition().X * d.PTM), float32(body.GetPosition().Y * d.PTM), 0}, float32(circle.M_radius*d.PTM), 10, false)
		c.SetAnchorToCenter()
		c.SetColor(color)
		body.SetUserData(c)
	case box2d.B2Shape_Type.E_polygon:
		b2Shape := fixture.GetShape().(*box2d.B2PolygonShape)
		numVertices := b2Shape.M_count
		var vertices []mgl32.Vec2
		for i := 0; i < numVertices; i++ {
			vertices = append(vertices, mgl32.Vec2{float32(b2Shape.M_vertices[i].X), float32(b2Shape.M_vertices[i].Y)})
		}
		c := graphics.NewPolylinePrimitive(mgl32.Vec3{float32(body.GetPosition().X * d.PTM), float32(body.GetPosition().Y * d.PTM), 0}, vertices, true)
		c.SetScale(mgl32.Vec2{float32(d.PTM), float32(d.PTM)})
		c.SetAnchorToCenter()
		c.SetColor(color)
		body.SetUserData(c)
	case box2d.B2Shape_Type.E_chain:
		b2Shape := fixture.GetShape().(*box2d.B2ChainShape)
		numVertices := b2Shape.M_count
		var vertices []mgl32.Vec2
		for i := 0; i < numVertices; i++ {
			vertices = append(vertices, mgl32.Vec2{float32(b2Shape.M_vertices[i].X), float32(b2Shape.M_vertices[i].Y)})
		}
		c := graphics.NewPolylinePrimitive(mgl32.Vec3{float32(body.GetPosition().X * d.PTM), float32(body.GetPosition().Y * d.PTM), 0}, vertices, false)
		c.SetScale(mgl32.Vec2{float32(d.PTM), float32(d.PTM)})
		c.SetAnchorToCenter()
		c.SetColor(color)
		body.SetUserData(c)

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
		c := body.GetUserData().(*graphics.Primitive2D)
		c.SetPosition(mgl32.Vec3{float32(body.GetPosition().X * d.PTM), float32(body.GetPosition().Y * d.PTM), 0})
		c.SetAngle(float32(body.GetAngle()))
		c.SetColor(color)
	case box2d.B2Shape_Type.E_polygon:
		c := body.GetUserData().(*graphics.Primitive2D)
		c.SetPosition(mgl32.Vec3{float32(body.GetPosition().X * d.PTM), float32(body.GetPosition().Y * d.PTM), 0})
		c.SetAngle(float32(body.GetAngle()))
		c.SetColor(color)
	case box2d.B2Shape_Type.E_chain:
		c := body.GetUserData().(*graphics.Primitive2D)
		c.SetPosition(mgl32.Vec3{float32(body.GetPosition().X * d.PTM), float32(body.GetPosition().Y * d.PTM), 0})
		c.SetAngle(float32(body.GetAngle()))
		c.SetColor(color)

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

// 	if (flags & b2Draw::e_jointBit)
// 	{
// 		for (b2Joint* j = m_jointList; j; j = j.GetNext())
// 		{
// 			DrawJoint(j);
// 		}
// 	}

// 	if (flags & b2Draw::e_pairBit)
// 	{
// 		b2Color color(0.3f, 0.9f, 0.9f);
// 		for (b2Contact* c = m_contactManager.m_contactList; c; c = c.GetNext())
// 		{
// 			//b2Fixture* fixtureA = c.GetFixtureA();
// 			//b2Fixture* fixtureB = c.GetFixtureB();

// 			//b2Vec2 cA = fixtureA.GetAABB().GetCenter();
// 			//b2Vec2 cB = fixtureB.GetAABB().GetCenter();

// 			//g_debugDraw.DrawSegment(cA, cB, color);
// 		}
// 	}

// 	if (flags & b2Draw::e_aabbBit)
// 	{
// 		b2Color color(0.9f, 0.3f, 0.9f);
// 		b2BroadPhase* bp = &m_contactManager.m_broadPhase;

// 		for (b2Body* b = m_bodyList; b; b = b.GetNext())
// 		{
// 			if (b.IsActive() == false)
// 			{
// 				continue;
// 			}

// 			for (b2Fixture* f = b.GetFixtureList(); f; f = f.GetNext())
// 			{
// 				for (int i = 0; i < f.m_proxyCount; ++i)
// 				{
// 					b2FixtureProxy* proxy = f.m_proxies + i;
// 					b2AABB aabb = bp.GetFatAABB(proxy.proxyId);
// 					b2Vec2 vs[4];
// 					vs[0].Set(aabb.lowerBound.x, aabb.lowerBound.y);
// 					vs[1].Set(aabb.upperBound.x, aabb.lowerBound.y);
// 					vs[2].Set(aabb.upperBound.x, aabb.upperBound.y);
// 					vs[3].Set(aabb.lowerBound.x, aabb.upperBound.y);

// 					g_debugDraw.DrawPolygon(vs, 4, color);
// 				}
// 			}
// 		}
// 	}

// 	if (flags & b2Draw::e_centerOfMassBit)
// 	{
// 		for (b2Body* b = m_bodyList; b; b = b.GetNext())
// 		{
// 			b2Transform xf = b.GetTransform();
// 			xf.p = b.GetWorldCenter();
// 			g_debugDraw.DrawTransform(xf);
// 		}
// 	}
// }

// void (world *B2World) DrawJoint(b2Joint* joint)
// {
// 	b2Body* bodyA = joint.GetBodyA();
// 	b2Body* bodyB = joint.GetBodyB();
// 	const b2Transform& xf1 = bodyA.GetTransform();
// 	const b2Transform& xf2 = bodyB.GetTransform();
// 	b2Vec2 x1 = xf1.p;
// 	b2Vec2 x2 = xf2.p;
// 	b2Vec2 p1 = joint.GetAnchorA();
// 	b2Vec2 p2 = joint.GetAnchorB();

// 	b2Color color(0.5f, 0.8f, 0.8f);

// 	switch (joint.GetType())
// 	{
// 	case e_distanceJoint:
// 		g_debugDraw.DrawSegment(p1, p2, color);
// 		break;

// 	case e_pulleyJoint:
// 		{
// 			b2PulleyJoint* pulley = (b2PulleyJoint*)joint;
// 			b2Vec2 s1 = pulley.GetGroundAnchorA();
// 			b2Vec2 s2 = pulley.GetGroundAnchorB();
// 			g_debugDraw.DrawSegment(s1, p1, color);
// 			g_debugDraw.DrawSegment(s2, p2, color);
// 			g_debugDraw.DrawSegment(s1, s2, color);
// 		}
// 		break;

// 	case e_mouseJoint:
// 		// don't draw this
// 		break;

// 	default:
// 		g_debugDraw.DrawSegment(x1, p1, color);
// 		g_debugDraw.DrawSegment(p1, p2, color);
// 		g_debugDraw.DrawSegment(x2, p2, color);
// 	}
// }
