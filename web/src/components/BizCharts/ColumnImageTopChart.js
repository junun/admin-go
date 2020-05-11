import React from "react";
import { Radio, Card, Col, Row, Spin} from "antd";

import styles from './ColumnImageTopChart.less'

import {
  G2,
  Chart,
  Geom,
  Axis,
  Tooltip,
  Coord,
  Label,
  Legend,
  View,
  Guide,
  Shape,
  Facet,
  Util
} from "bizcharts";


class ColumnImageTopChart extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      chart: null,
    };
  }

  onGetG2Instance = (g2Chart) => {
    this.setState({
      chart: g2Chart,
    });
  };

  onPlotClick = (ev) => {
    const point = {
      x: ev.x,
      y: ev.y
    };
    const items = this.state.chart.getTooltipItems(point);
    if (items && items.length > 0) {
      this.props.onPlotClick(items[0]);
    }
  };

  render() {
    const { loading, data, title, onDateChange, type } = this.props;

    let chart;
    if (data && data.length > 0) {
      chart = <Chart
        onGetG2Instance={this.onGetG2Instance}
        onPlotClick={this.onPlotClick}
        forceFit={true}
        height={400}
        data={data}
        style={{ marginTop: '5px'}}>
          <Axis name="Name" />
          <Axis name="Amount" />
          <Tooltip />
          <Geom type="interval" position="Name*Amount" />
        </Chart>;
    } else {
      chart = <Chart placeholder height={400} />;
    }

    return <Card title={title} style={{height: '500px'}}>
      <Row style={{ marginTop: '5px'}}>
        <Col span={16}>
          <Radio.Group value={type} onChange={onDateChange}>
            <Radio.Button value="0">昨日</Radio.Button>
            <Radio.Button value="1">周</Radio.Button>
            <Radio.Button value="2">月</Radio.Button>
            <Radio.Button value="3">季度</Radio.Button>
            <Radio.Button value="4">年度</Radio.Button>
          </Radio.Group>
        </Col>
      </Row>
      {loading ?
        <Spin><div style={{ height: '400px'}}></div></Spin>
      :
        chart
      }
    </Card>;
  }
}

export default ColumnImageTopChart;

