import React from 'react';
import { Radio, Card, Col, Row, Spin} from "antd"
import {
  Chart,
  Geom,
  Tooltip,
  Coord,
  Label,
  View,
} from 'bizcharts';
import DataSet from '@antv/data-set';

class SunburstPieChart extends React.Component {
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
    const { loading, title, data,  onDateChange, type } = this.props;

    const { DataView } = DataSet;
    const dv = new DataView();
    dv.source(data).transform({
      type: 'percent',
      field: 'Total',
      dimension: 'Name',
      as: 'percent',
    });
    
    const cols = {
      percent: {
        formatter: (val) => {
          val = `${(val * 100).toFixed(2)}%`;
          return val;
        },
      },
    };
    
    const dv1 = new DataView();
    dv1.source(data).transform({
      type: 'percent',
      field: 'Total',
      dimension: 'PackageName',
      as: 'percent',
    });
    
    let chart;
    if (data && data.length > 0) {
      chart = <Chart
        onGetG2Instance={this.onGetG2Instance}
        onPlotClick={this.onPlotClick}
        forceFit={true}
        height={400}
        data={dv}
        scale={cols}
        style={{ marginTop: '5px'}}>
          <Coord type="theta" radius={0.5} />
          <Tooltip
            showTitle={false}
            itemTpl="<li>
                      <span style=&quot;background-color:{color};&quot; class=&quot;g2-tooltip-marker&quot;>
                      </span>{name}: {value}
                    </li>"
          />
          <Geom
            type="intervalStack"
            position="percent"
            color="Name"
            tooltip={[
              'Name*percent',
              (item, percent) => {
                percent = `${(percent * 100).toFixed(2)}%`;
                return {
                  name: item,
                  value: percent,
                };
              },
            ]}
            style={{
              lineWidth: 1,
              stroke: '#fff',
            }}
            select={false}
          >
            <Label content="Name" offset={-10} />
          </Geom>
          <View data={dv1} scale={cols}>
            <Coord type="theta" radius={0.75} innerRadius={0.5 / 0.75} />
            <Geom
              type="intervalStack"
              position="percent"
              color={[
                'PackageName',
                [
                  '#BAE7FF',
                  '#7FC9FE',
                  '#71E3E3',
                  '#ABF5F5',
                  '#8EE0A1',
                  '#BAF5C4',
                ],
              ]}
              tooltip={[
                'PackageName*percent',
                (PackageName, percent) => {
                  percent = `${(percent * 100).toFixed(2)}%`;
                  return {
                    name: PackageName,
                    value: percent,
                  };
                },
              ]}
              style={{
                lineWidth: 1,
                stroke: '#fff',
              }}
              select={false}
            >
              <Label content="PackageName" />
            </Geom>
          </View>
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

export default SunburstPieChart;
