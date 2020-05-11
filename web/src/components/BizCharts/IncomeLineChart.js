import React, {PureComponent} from "react";
import { Radio, Card, Col, Row, Spin, DatePicker} from "antd";
import { Chart, Geom, Axis, Tooltip } from 'bizcharts';
import styles from './IncomeLineChart.less'

const { RangePicker } = DatePicker;

class IncomeLineChart extends PureComponent {
  constructor(props) {
    super(props);
    this.state = {
      chart: null,
    };
  }

  componentDidMount() {
    const e = document.createEvent("Event");
    e.initEvent("resize", true, true);
    window.dispatchEvent(e);
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
    const { loading, title, data, onDateChange, onDatePick, beginDate, endDate, type } = this.props;
    let chart;
    if (data && data.length > 0) {
      const amountLabelConfig = {
        formatter(text, item, index) {
          const amount = parseInt(text);
          return (amount / 100).toString();
        },
      };
      chart = <Chart onGetG2Instance={this.onGetG2Instance} onPlotClick={this.onPlotClick} forceFit={true} height={300} data={data} style={{ marginTop: '5px'}} padding={[20, 40, 30, 80]}>
        <Axis name="time"/>
        <Axis name="amount" label={amountLabelConfig} />
        <Tooltip />
        <Geom type="line" position="time*amount" color="#97c982" tooltip={['time*amount*quantity', (time, amount) => {
          return {
            name: '点灯总金额',
            title: time,
            value: (amount / 100).toFixed(2),
          };
        }]} />
      </Chart>;
    } else {
      chart = <Chart placeholder height={300} />;
    }
    let sum = 0;
    if (data && data.length > 0) {
      sum = data.map(x => x.amount).reduce((total, amount )=> total + amount);
    }
    return <Card title={title} style={{height: '492px'}}>
      <Row>
        <Col>
          <span className={styles.chartSumLabel}>{(sum / 100).toFixed(2)}</span>
          <span className={styles.chartSumTitle}>元</span>
        </Col>
      </Row>
      <Row style={{ marginTop: '5px'}}>
        <Col span={10}>
          <Radio.Group value={type} onChange={onDateChange}>
            <Radio.Button value="0">今日</Radio.Button>
            <Radio.Button value="1">周</Radio.Button>
            <Radio.Button value="2">月</Radio.Button>
            <Radio.Button value="3">季度</Radio.Button>
            <Radio.Button value="4">年度</Radio.Button>
          </Radio.Group>
        </Col>
        <Col span={14}>
          <RangePicker value={[beginDate, endDate]} onChange={onDatePick} style={{ width: '280px', float: 'right' }} />
        </Col>
      </Row>
      {loading ?
        <Spin><div style={{ height: '300px'}}></div></Spin>
      :
        chart
      }
    </Card>;
  }
}

export default IncomeLineChart;
