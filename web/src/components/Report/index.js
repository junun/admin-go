import React, {PureComponent} from "react";
import { Card } from 'antd';

class ReportCard extends PureComponent {
  render() {
    const { reports } = this.props;
    return <Card title="Message For You">
      <pre>{reports}</pre>
    </Card>;
  }
}

export default ReportCard;